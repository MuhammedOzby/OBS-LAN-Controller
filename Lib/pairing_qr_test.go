package lib

import (
	"bytes"
	"encoding/base64"
	"html"
	"image/png"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func TestPairingQRDecodeAndBoundary(t *testing.T) {
	router, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	h := newPairingHandler(router)
	h.listenAddr = &net.TCPAddr{IP: net.IPv4zero, Port: 43217}
	v, err := h.view(true, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(v.QRs) == 0 {
		t.Skip("no non-loopback IPv4 interface available")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(string(v.QRs[0].Image), "data:image/png;base64,"))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := qrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(decoded.GetText())
	if err != nil {
		t.Fatal(err)
	}
	if u.Host != v.QRs[0].Address || u.Path != "/pair" || u.RawQuery != "" || u.Fragment != v.Code {
		t.Fatal("decoded QR destination incorrect")
	}
	remote := "192.0.2.3:1234"
	r := pairingRequest(h, "GET", "/pair", remote, "", "", nil)
	if r.Code != 200 || strings.Contains(r.Body.String(), v.Code) || strings.Contains(r.Body.String(), "data:image/png") {
		t.Fatal("unpaired client received secret")
	}
	post := pairingRequest(h, "POST", "/pair", remote, "http://localhost:8080", url.Values{"code": {u.Fragment}}.Encode(), nil)
	if post.Code != 303 {
		t.Fatalf("redemption status %d", post.Code)
	}
	cookie := post.Result().Cookies()[0]
	if pairingRequest(h, "POST", "/pair", remote, "http://localhost:8080", "code="+v.Code, nil).Code != 401 {
		t.Fatal("reused QR accepted")
	}
	local := pairingRequest(h, "GET", "/mobil-cihaz-ile-giris", "127.0.0.1:1", "", "", nil)
	if local.Code != 200 || !strings.Contains(local.Body.String(), "data:image/png;base64,") || strings.Count(local.Body.String(), `class="nav-link`) != 4 || strings.Contains(local.Body.String(), "data-broadcast-action") {
		t.Fatal("local mobile tab incorrect")
	}
	code := h.code
	paired := pairingRequest(h, "GET", "/mobil-cihaz-ile-giris", remote, "", "", cookie)
	if paired.Code != 200 || strings.Contains(paired.Body.String(), code) || strings.Contains(paired.Body.String(), "data:image/png") || strings.Contains(paired.Body.String(), "revoke-all") || h.code != code {
		t.Fatal("paired remote accessed issuance")
	}
	if pairingRequest(h, "GET", "/mobil-cihaz-ile-giris", remote, "", "", nil).Code != 303 {
		t.Fatal("unpaired tab accessible")
	}
	if pairingRequest(h, "POST", "/pair", remote, "http://evil.test", "code="+code, nil).Code != 403 {
		t.Fatal("cross-origin accepted")
	}
	h.listenAddr = &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 43217}
	loopback := pairingRequest(h, "GET", "/pair", "127.0.0.1:1", "", "", nil)
	if strings.Contains(loopback.Body.String(), "data:image/png") || !strings.Contains(loopback.Body.String(), "yeniden başlatın") {
		t.Fatal("loopback-only instructions missing")
	}
	server, err := newHTTPServer(":43217", h)
	if err != nil {
		t.Fatal(err)
	}
	if pairingRequest(server.Handler, "GET", "/pair", "127.0.0.1:1", "http://evil.test", "", nil).Code != http.StatusForbidden {
		t.Fatal("spoofed origin accepted")
	}
}

func TestPairingQRChromeFlow(t *testing.T) {
	chrome := os.Getenv("CHROME_BIN")
	if chrome == "" {
		t.Skip("CHROME_BIN not set")
	}
	listener, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("<html><body>CONTROL_OK</body></html>")) }))
	h.listenAddr = listener.Addr()
	if len(pairingAddresses(h.listenAddr)) == 0 {
		t.Skip("no LAN interface")
	}
	server, err := newHTTPServer(listener.Addr().String(), h)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Serve(listener)
	local, err := http.Get("http://127.0.0.1:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port) + "/pair")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(local.Body)
	local.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	match := regexp.MustCompile(`src="data:image/png;base64,([^"]+)"`).FindSubmatch(body)
	if len(match) != 2 {
		t.Fatal("local QR image missing")
	}
	data, err := base64.StdEncoding.DecodeString(html.UnescapeString(string(match[1])))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		t.Fatal(err)
	}
	result, err := qrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage", "--no-proxy-server", "--user-data-dir="+t.TempDir(), "--virtual-time-budget=3000", "--dump-dom", result.GetText())
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Chrome failed: %v", err)
	}
	if !strings.Contains(string(out), "CONTROL_OK") {
		t.Fatal("QR navigation did not reach control")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.sessions) != 1 || h.code != "" {
		t.Fatal("Chrome did not consume code and create session")
	}
}
