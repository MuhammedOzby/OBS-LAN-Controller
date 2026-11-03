package lib

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPServerTrustedHosts(t *testing.T) {
	called := false
	server, err := newHTTPServer(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		host string
		want int
	}{
		{"localhost", 204}, {"LOCALHOST:8080", 204}, {"localhost.:8080", 204},
		{"127.0.0.1:8080", 204}, {"127.20.30.40", 204}, {"[::1]:8080", 204},
		{"[::ffff:127.0.0.1]:8080", 204},
		{"attacker.example:8080", 403}, {"controller.test", 403},
		{"localhost.attacker.example", 403}, {"192.0.2.123:8080", 403},
		{"[2001:db8::123]:8080", 403}, {"0.0.0.0:8080", 403}, {"[::]:8080", 403},
		{"127.0.0.1:invalid", 403}, {"localhost:65536", 403},
		{"localhost:", 403}, {"localhost@attacker.example", 403},
		{"", 403}, {"::1", 403},
	} {
		t.Run(tc.host, func(t *testing.T) {
			called = false
			r := httptest.NewRequest("GET", "http://localhost/api/status", nil)
			r.Host = tc.host
			// Neither forwarding headers nor same-origin attacker input may bypass Host validation.
			r.Header.Set("X-Forwarded-Host", "localhost:8080")
			r.Header.Set("Forwarded", "host=localhost:8080")
			r.Header.Set("Origin", "http://"+tc.host)
			rr := httptest.NewRecorder()
			server.Handler.ServeHTTP(rr, r)
			if rr.Code != tc.want || called != (tc.want == 204) {
				t.Fatalf("Host %q: status=%d called=%v, want %d", tc.host, rr.Code, called, tc.want)
			}
		})
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		t.Fatal(err)
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil || ip.IsUnspecified() {
			continue
		}
		r := httptest.NewRequest("GET", "http://localhost/", nil)
		r.Host = net.JoinHostPort(ip.String(), "8080")
		rr := httptest.NewRecorder()
		server.Handler.ServeHTTP(rr, r)
		if rr.Code != 204 {
			t.Errorf("assigned interface address %q rejected: %d", r.Host, rr.Code)
		}
	}
}

// Real TCP connections send valid headers but only one byte of the promised
// JSON body. A context timeout alone cannot unblock these socket reads.
func TestHTTPServerBoundsSlowActionBodies(t *testing.T) {
	for _, path := range []string{"/api/scene", "/api/broadcast", "/api/deck/action"} {
		t.Run(path, func(t *testing.T) {
			h, err := newHandler(nil)
			if err != nil {
				t.Fatal(err)
			}
			finished := make(chan struct{}, 1)
			server, err := newHTTPServer("127.0.0.1:0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() { finished <- struct{}{} }()
				h.ServeHTTP(w, r)
			}))
			if err != nil {
				t.Fatal(err)
			}
			if server.ReadTimeout <= 0 || server.ReadTimeout > 10*time.Second {
				t.Errorf("production body ReadTimeout = %s, want positive and <=10s", server.ReadTimeout)
			} else {
				// Same production server, shortened deadline for deterministic tests.
				server.ReadTimeout = 100 * time.Millisecond
			}
			listener, err := net.Listen("tcp", server.Addr)
			if err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- server.Serve(listener) }()
			t.Cleanup(func() { server.Close(); <-done })
			conn, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			host := listener.Addr().String()
			if _, err := fmt.Fprintf(conn, "POST %s HTTP/1.1\r\nHost: %s\r\nOrigin: http://%s\r\nContent-Type: application/json\r\nContent-Length: 100\r\n\r\n{", path, host, host); err != nil {
				t.Fatal(err)
			}
			select {
			case <-finished:
			case <-time.After(800 * time.Millisecond):
				t.Fatal("action handler remained blocked reading incomplete body")
			}
			if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
				t.Fatal(err)
			}
			response, err := http.ReadResponse(bufio.NewReader(conn), nil)
			if err != nil {
				t.Fatalf("read timeout response: %v", err)
			}
			defer response.Body.Close()
			io.Copy(io.Discard, response.Body)
			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("incomplete body status = %d, want 400", response.StatusCode)
			}
		})
	}
}

func TestHTTPServerGuardPreservesOriginChecksAndFixtures(t *testing.T) {
	h, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	server, err := newHTTPServer(":8080", h)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/scene", "/api/broadcast", "/api/deck/action"} {
		for _, tc := range []struct {
			host, origin string
			want         int
		}{
			{"attacker.example", "http://attacker.example", 403},
			{"127.0.0.1:8080", "http://attacker.example", 403},
			{"127.0.0.1:8080", "", 403},
		} {
			r := httptest.NewRequest("POST", "http://localhost"+path, strings.NewReader(`{}`))
			r.Host = tc.host
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			server.Handler.ServeHTTP(rr, r)
			if rr.Code != tc.want {
				t.Errorf("%s Host=%q Origin=%q: %d, want %d", path, tc.host, tc.origin, rr.Code, tc.want)
			}
		}
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest("GET", "http://controller.test/", nil))
	if rr.Code != 200 {
		t.Fatalf("unguarded router fixture status=%d", rr.Code)
	}
}
