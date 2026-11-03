//go:build linux

package lib

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"path/filepath"
)

func ownedPreviewPath(path string, mode os.FileMode) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	var stat unix.Stat_t
	if err := unix.Lstat(path, &stat); err != nil {
		return err
	}
	if stat.Uid != uint32(os.Geteuid()) || info.Mode().Perm()&0077 != 0 || info.Mode()&os.ModeType != mode {
		return fmt.Errorf("unsafe preview path: %s (same owner and private permissions required)", path)
	}
	return nil
}

func listenPreview() (*net.UnixListener, func(), error) {
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if !filepath.IsAbs(runtime) {
		return nil, nil, fmt.Errorf("preview requires XDG_RUNTIME_DIR; run Controller as the OBS desktop user")
	}
	if err := ownedPreviewPath(runtime, os.ModeDir); err != nil {
		return nil, nil, err
	}
	dir := filepath.Join(runtime, "obs-lan-controller")
	if err := os.Mkdir(dir, 0700); err != nil && !os.IsExist(err) {
		return nil, nil, err
	}
	if err := ownedPreviewPath(dir, os.ModeDir); err != nil {
		return nil, nil, err
	}
	lockPath := filepath.Join(dir, "preview.lock")
	fd, err := unix.Open(lockPath, unix.O_RDWR|unix.O_CREAT|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0600)
	if err != nil {
		return nil, nil, err
	}
	unlock := func() { _ = unix.Close(fd) } // keep inode: never unlink an advisory lock
	var stat unix.Stat_t
	if err = unix.Fstat(fd, &stat); err != nil || stat.Uid != uint32(os.Geteuid()) || stat.Mode&unix.S_IFMT != unix.S_IFREG || stat.Mode&0077 != 0 || stat.Nlink != 1 {
		unlock()
		return nil, nil, fmt.Errorf("unsafe preview lock")
	}
	if err := unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); err != nil {
		unlock()
		return nil, nil, fmt.Errorf("preview already running: %w", err)
	}
	path := filepath.Join(dir, "preview.sock")
	if len(path) >= 108 {
		unlock()
		return nil, nil, fmt.Errorf("preview socket path too long")
	}
	if _, err := os.Lstat(path); err == nil {
		if err = ownedPreviewPath(path, os.ModeSocket); err != nil {
			unlock()
			return nil, nil, err
		}
		conn, dialErr := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
		if dialErr == nil {
			conn.Close()
			unlock()
			return nil, nil, fmt.Errorf("preview socket is already in use")
		}
		if !errors.Is(dialErr, unix.ECONNREFUSED) {
			unlock()
			return nil, nil, dialErr
		}
		if err := os.Remove(path); err != nil {
			unlock()
			return nil, nil, err
		}
	} else if !os.IsNotExist(err) {
		unlock()
		return nil, nil, err
	}
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		unlock()
		return nil, nil, err
	}
	if err := os.Chmod(path, 0600); err != nil {
		listener.Close()
		unlock()
		return nil, nil, err
	}
	return listener, func() { listener.Close(); unlock() }, nil
}

func previewPeerAllowed(conn *net.UnixConn) bool {
	raw, err := conn.SyscallConn()
	if err != nil {
		return false
	}
	allowed := false
	err = raw.Control(func(fd uintptr) {
		cred, err := unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
		allowed = err == nil && cred.Uid == uint32(os.Geteuid())
	})
	return err == nil && allowed
}
