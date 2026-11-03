//go:build !linux

package lib

import (
	"fmt"
	"net"
)

func listenPreview() (*net.UnixListener, func(), error) {
	return nil, nil, fmt.Errorf("native preview socket is currently supported only on Linux")
}
func previewPeerAllowed(*net.UnixConn) bool { return false }
