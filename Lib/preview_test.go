//go:build linux

package lib

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"
)

func testRecord(kind byte, generation uint32, key bool, pts int64, data []byte) []byte {
	out := new(bytes.Buffer)
	_ = binary.Write(out, binary.BigEndian, uint32(36+len(data)))
	out.WriteByte(kind)
	if key {
		out.WriteByte(1)
	} else {
		out.WriteByte(0)
	}
	_ = binary.Write(out, binary.BigEndian, uint16(0))
	_ = binary.Write(out, binary.BigEndian, generation)
	_ = binary.Write(out, binary.BigEndian, pts)
	_ = binary.Write(out, binary.BigEndian, pts)
	_ = binary.Write(out, binary.BigEndian, uint32(1))
	_ = binary.Write(out, binary.BigEndian, uint32(30))
	_ = binary.Write(out, binary.BigEndian, uint32(len(data)))
	out.Write(data)
	return out.Bytes()
}
func testConfig() []byte {
	b := new(bytes.Buffer)
	for _, v := range []uint32{1280, 720, 30, 1} {
		_ = binary.Write(b, binary.BigEndian, v)
	}
	b.Write([]byte{0, 0, 0, 1, 0x67, 0x42, 0xe0, 0x1f, 0, 0, 0, 1, 0x68, 0xce, 0x06})
	return b.Bytes()
}
func TestPreviewProtocolBounds(t *testing.T) {
	for _, raw := range [][]byte{{0, 0, 0, 1}, {0xff, 0xff, 0xff, 0xff}, testRecord(4, 1, false, 0, []byte{1}), testRecord(3, 0, false, 0, []byte{1})} {
		if _, err := readPreviewRecord(bytes.NewReader(raw)); err == nil {
			t.Fatal("invalid record accepted")
		}
	}
	record, err := readPreviewRecord(bytes.NewReader(testRecord(2, 1, false, 0, testConfig())))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parsePreviewConfig(record); err != nil {
		t.Fatal(err)
	}
	p := newPreviewService()
	defer p.close()
	if err := p.consume(record); err != nil {
		t.Fatal(err)
	}
	frame, _ := readPreviewRecord(bytes.NewReader(testRecord(3, 1, true, 1, []byte{0, 0, 1, 0x65, 0x88})))
	if err := p.consume(frame); err != nil {
		t.Fatal(err)
	}
	if err := p.consume(frame); err == nil {
		t.Fatal("non-monotonic PTS accepted")
	}
	frame.generation = 2
	if err := p.consume(frame); err == nil {
		t.Fatal("wrong generation accepted")
	}
	record.data[0] = 0xff
	if _, err := parsePreviewConfig(record); err == nil {
		t.Fatal("oversize width accepted")
	}
}
func setupPreview(t *testing.T) (*previewService, string) {
	t.Helper()
	// Short Unix socket path under the designated scratch root.
	runtime, err := os.MkdirTemp(os.Getenv("TMPDIR"), "pv-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(runtime) })
	t.Setenv("XDG_RUNTIME_DIR", runtime)
	p := newPreviewService()
	if err := p.start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(p.close)
	return p, filepath.Join(runtime, "obs-lan-controller/preview.sock")
}
func connectTestSource(t *testing.T, path string) *net.UnixConn {
	t.Helper()
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte(previewMagic)); err != nil {
		t.Fatal(err)
	}
	var hello [8]byte
	if _, err := io.ReadFull(conn, hello[:]); err != nil || string(hello[:]) != previewMagic {
		t.Fatalf("handshake: %v %q", err, hello)
	}
	return conn
}
func awaitCommand(t *testing.T, conn *net.UnixConn, want byte) {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))
	var b [1]byte
	for {
		if _, err := io.ReadFull(conn, b[:]); err != nil {
			t.Fatalf("waiting %c: %v", want, err)
		}
		if b[0] == want {
			return
		}
	}
}
func TestPreviewSocketLifecycle(t *testing.T) {
	p, path := setupPreview(t)
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatal("socket not private", err)
	}
	second := newPreviewService()
	if err := second.start(); err == nil {
		second.close()
		t.Fatal("second Controller acquired socket")
	}
	second.close()
	conn := connectTestSource(t, path)
	awaitCommand(t, conn, 'P')
	p.mu.Lock()
	p.viewers["test"] = &previewViewer{ready: true}
	p.mu.Unlock()
	p.notify()
	awaitCommand(t, conn, 'S')
	conn.Write(testRecord(2, 1, false, 0, testConfig()))
	p.mu.Lock()
	delete(p.viewers, "test")
	p.mu.Unlock()
	p.notify()
	awaitCommand(t, conn, 'T')
	conn.Close()
	deadline := time.Now().Add(time.Second)
	for p.status().Connected && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if p.status().Connected {
		t.Fatal("disconnect not reflected")
	}
	conn = connectTestSource(t, path)
	awaitCommand(t, conn, 'P')
}
func TestPreviewUnsafeRuntime(t *testing.T) {
	runtime := t.TempDir()
	os.Chmod(runtime, 0755)
	t.Setenv("XDG_RUNTIME_DIR", runtime)
	p := newPreviewService()
	defer p.close()
	if err := p.start(); err == nil {
		t.Fatal("public runtime accepted")
	}
}
func TestPreviewOfferSecurity(t *testing.T) {
	p := newPreviewService()
	defer p.close()
	handler, err := newHandlerWithPreview(nil, p)
	if err != nil {
		t.Fatal(err)
	}
	handler = newPairingHandler(handler)
	for _, tc := range []struct {
		peer, origin string
		code         int
	}{{"192.0.2.2:1234", "http://example.test", 401}, {"127.0.0.1:1234", "http://evil.test", 403}, {"127.0.0.1:1234", "http://example.test", 400}} {
		req := httptest.NewRequest("POST", "http://example.test/api/preview/offer", strings.NewReader(`{"type":"answer","sdp":"invalid"}`))
		req.RemoteAddr = tc.peer
		req.Header.Set("Origin", tc.origin)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != tc.code {
			t.Fatalf("got %d want %d", rec.Code, tc.code)
		}
	}
}
func TestPreviewOfferLimit(t *testing.T) {
	p := newPreviewService()
	defer p.close()
	h, _ := newHandlerWithPreview(nil, p)
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	defer pc.Close()
	for i := 0; i < previewMaxViewers; i++ {
		p.viewers[string(rune('a'+i))] = &previewViewer{id: string(rune('a' + i)), pc: pc, done: make(chan struct{})}
	}
	req := httptest.NewRequest(http.MethodPost, "http://example.test/api/preview/offer", strings.NewReader(`{"type":"offer","sdp":"unused"}`))
	req.Header.Set("Origin", "http://example.test")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 503 {
		t.Fatal(rec.Code)
	}
}
