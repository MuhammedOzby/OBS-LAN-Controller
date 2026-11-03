package lib

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

func pairingRequest(h http.Handler, method, path, remote, origin, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://localhost:8080"+path, strings.NewReader(body))
	r.RemoteAddr = remote
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	if method == "POST" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, r)
	return rr
}

func pairingCode(t *testing.T, h http.Handler) string {
	t.Helper()
	rr := pairingRequest(h, "GET", "/pair", "127.0.0.1:5000", "", "", nil)
	if rr.Code != 200 {
		t.Fatalf("local pair status = %d", rr.Code)
	}
	matches := regexp.MustCompile(`<output[^>]*>([A-F0-9]{12})</output>`).FindStringSubmatch(rr.Body.String())
	if len(matches) != 2 {
		t.Fatal("local pairing code missing")
	}
	if rr.Header().Get("Cache-Control") != "no-store" {
		t.Error("pair page must not be cached")
	}
	return matches[1]
}

func TestPairingLifecycle(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("authorized readback")) }))
	code := pairingCode(t, h)
	remote := pairingRequest(h, "GET", "/pair", "192.0.2.3:5000", "", "", nil)
	if remote.Code != 200 || strings.Contains(remote.Body.String(), code) || !strings.Contains(remote.Body.String(), `action="/pair"`) {
		t.Fatal("remote form must not disclose code")
	}
	paired := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", url.Values{"code": {code}}.Encode(), nil)
	if paired.Code != 303 || paired.Header().Get("Location") != "/" {
		t.Fatalf("pairing status %d", paired.Code)
	}
	cookies := paired.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatal("missing session")
	}
	c := cookies[0]
	if !c.HttpOnly || c.SameSite != http.SameSiteStrictMode || c.Path != "/" || c.MaxAge <= 0 || c.Secure || len(c.Value) < 40 || c.Domain != "" {
		t.Fatal("unsafe HTTP cookie attributes")
	}
	if c.Expires.IsZero() {
		t.Fatal("missing cookie expiry")
	}
	readback := pairingRequest(h, "GET", "/api/status", "192.0.2.3:5000", "", "", c)
	if readback.Code != 200 || readback.Body.String() != "authorized readback" {
		t.Fatal("session readback failed")
	}
	reused := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil)
	if reused.Code != 401 {
		t.Fatalf("single use got %d", reused.Code)
	}
	forged := *c
	forged.Value += "x"
	if pairingRequest(h, "GET", "/api/status", "192.0.2.3:5000", "", "", &forged).Code != 401 {
		t.Fatal("forged session accepted")
	}
}

func TestPairingExpiryRotationAndLimits(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h.now = func() time.Time { return now }
	post := func(code string) *httptest.ResponseRecorder {
		return pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil)
	}
	code := pairingCode(t, h)
	next := pairingCode(t, h)
	if next == code {
		t.Fatal("refresh must generate a new code")
	}
	if post(code).Code != 401 || post("wrong").Code != 401 {
		t.Fatal("old or wrong code accepted")
	}
	now = now.Add(pairingCodeTTL)
	if post(next).Code != 401 {
		t.Fatal("expired code accepted at boundary")
	}
	code = pairingCode(t, h)
	paired := post(code)
	if paired.Code != 303 {
		t.Fatal("pairing failed")
	}
	cookie := paired.Result().Cookies()[0]
	now = now.Add(pairingSessionTTL)
	if pairingRequest(h, "GET", "/api/status", "192.0.2.3:5000", "", "", cookie).Code != 401 {
		t.Fatal("expired session accepted")
	}
	if len(h.sessions) != 0 {
		t.Fatal("expired sessions not pruned")
	}
	code = pairingCode(t, h)
	for i := 0; i < pairingMaxFailures; i++ {
		if post("wrong").Code != 401 {
			t.Fatal("premature lockout")
		}
	}
	for i := 0; i < 3; i++ {
		// New source addresses and refreshing the local page cannot reset the global limit.
		code = pairingCode(t, h)
		rr := pairingRequest(h, "POST", "/pair", "192.0.2.99:1234", "http://localhost:8080", "code="+code, nil)
		if rr.Code != 429 || rr.Header().Get("Retry-After") == "" {
			t.Fatal("global attempt limit missing")
		}
	}
	now = now.Add(pairingFailureWindow)
	if post(code).Code != 303 {
		t.Fatal("limit did not recover")
	}
	for len(h.sessions) < pairingMaxSessions {
		code = pairingCode(t, h)
		if post(code).Code != 303 {
			t.Fatal("session allocation failed")
		}
	}
	code = pairingCode(t, h)
	if post(code).Code != 503 || len(h.sessions) != pairingMaxSessions {
		t.Fatal("session storage unbounded")
	}
	now = now.Add(pairingSessionTTL)
	code = pairingCode(t, h)
	if post(code).Code != 303 || len(h.sessions) != 1 {
		t.Fatal("expired capacity not reclaimed")
	}
}

func TestPairingOriginsBodyAndTLS(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	code := pairingCode(t, h)
	for _, origin := range []string{"", "null", "http://evil.test", "https://localhost:8080", "http://localhost:8080/", "http://localhost:8080?x=y", "http://user@localhost:8080", "http://localhost:8080#x"} {
		rr := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", origin, "code="+code, nil)
		if rr.Code != 403 || len(rr.Result().Cookies()) != 0 {
			t.Errorf("forged Origin %q accepted: %d", origin, rr.Code)
		}
	}
	if pairingRequest(h, "GET", "/pair", "127.0.0.1:5000", "http://evil.test", "", nil).Code != 403 {
		t.Fatal("forged local issuance allowed")
	}
	r := httptest.NewRequest("GET", "http://localhost:8080/pair", nil)
	r.RemoteAddr = "127.0.0.1:5000"
	r.Header.Set("Sec-Fetch-Site", "cross-site")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, r)
	if rr.Code != 403 {
		t.Fatal("cross-site local issuance allowed")
	}
	if pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code+strings.Repeat("x", 2048), nil).Code != 400 {
		t.Fatal("oversized request accepted")
	}
	if pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code+"&code="+code, nil).Code != 401 {
		t.Fatal("duplicate form codes accepted")
	}
	if pairingRequest(h, "POST", "/pair?code="+code, "192.0.2.3:5000", "http://localhost:8080", "", nil).Code != 401 {
		t.Fatal("query-string code accepted")
	}
	if pairingRequest(h, "PUT", "/pair", "192.0.2.3:5000", "", "", nil).Code != 405 {
		t.Fatal("unexpected pair method allowed")
	}
	r = httptest.NewRequest("POST", "https://localhost:8080/pair", strings.NewReader("code="+code))
	r.RemoteAddr = "192.0.2.3:5000"
	r.TLS = &tls.ConnectionState{}
	r.Header.Set("Origin", "https://localhost:8080")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, r)
	if rr.Code != 303 || len(rr.Result().Cookies()) != 1 || !rr.Result().Cookies()[0].Secure {
		t.Fatal("TLS pairing must set Secure cookie")
	}
}

func TestPairingConcurrentSingleUse(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	code := pairingCode(t, h)
	results := make(chan int, 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil).Code
		}()
	}
	wg.Wait()
	close(results)
	successes := 0
	for status := range results {
		if status == 303 {
			successes++
		} else if status != 401 {
			t.Errorf("unexpected concurrent status %d", status)
		}
	}
	if successes != 1 {
		t.Fatalf("single-use successes = %d", successes)
	}
}

func TestPairingPreservesDownstreamOriginAndHostGuard(t *testing.T) {
	router, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	h := newPairingHandler(router)
	code := pairingCode(t, h)
	paired := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil)
	cookie := paired.Result().Cookies()[0]
	for _, path := range []string{"/api/scene", "/api/broadcast", "/api/deck/action"} {
		if pairingRequest(h, "POST", path, "192.0.2.3:5000", "http://evil.test", "{}", cookie).Code != 403 {
			t.Errorf("paired Origin guard bypass on %s", path)
		}
	}
	server, err := newHTTPServer(":8080", h)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "http://evil.test/pair", nil)
	r.RemoteAddr = "127.0.0.1:5000"
	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, r)
	if rr.Code != 403 || h.code != "" {
		t.Fatal("Host guard must precede pairing issuance")
	}
}

func TestPairingLocalRevokeAll(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	code := pairingCode(t, h)
	paired := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil)
	cookie := paired.Result().Cookies()[0]
	code = pairingCode(t, h)
	// LAN peers (including paired ones) cannot invoke local session revocation.
	if rr := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "action=revoke-all", cookie); rr.Code != 403 {
		t.Fatalf("remote revoke got %d", rr.Code)
	}
	if rr := pairingRequest(h, "POST", "/pair", "127.0.0.1:5000", "http://evil.test", "action=revoke-all", nil); rr.Code != 403 {
		t.Fatalf("forged revoke got %d", rr.Code)
	}
	if pairingRequest(h, "GET", "/api/status", "192.0.2.3:5000", "", "", cookie).Code != 204 {
		t.Fatal("unauthorized revoke invalidated session")
	}
	rr := pairingRequest(h, "POST", "/pair", "127.0.0.1:5000", "http://localhost:8080", "action=revoke-all", nil)
	if rr.Code != 303 || len(h.sessions) != 0 || h.code != "" {
		t.Fatalf("local revoke failed: status %d", rr.Code)
	}
	if pairingRequest(h, "GET", "/api/status", "192.0.2.3:5000", "", "", cookie).Code != 401 {
		t.Fatal("revoked session still authorized")
	}
	if pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", "code="+code, nil).Code != 401 {
		t.Fatal("revoke retained pending code")
	}
}

func TestPairingRemoteAccessBoundary(t *testing.T) {
	h := newPairingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	for _, tc := range []struct {
		path, remote string
		want         int
	}{
		{"/api/status", "192.0.2.4:5000", 401},
		{"/api/deck/action", "192.0.2.4:5000", 401},
		{"/deck", "192.0.2.4:5000", 303},
		{"/", "127.0.0.1:5000", 204},
		{"/api/status", "[::1]:5000", 204},
		{"/api/status", "[::ffff:127.0.0.1]:5000", 204},
		{"/api/status", "localhost:5000", 401},
		{"/api/status", "invalid", 401},
	} {
		r := httptest.NewRequest("GET", "http://localhost:8080"+tc.path, nil)
		r.RemoteAddr = tc.remote
		r.Header.Set("Forwarded", "for=127.0.0.1;host=localhost:8080")
		r.Header.Set("X-Forwarded-For", "127.0.0.1")
		r.Header.Set("X-Real-IP", "127.0.0.1")
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, r)
		if rr.Code != tc.want {
			t.Errorf("%s %s: got %d want %d", tc.path, tc.remote, rr.Code, tc.want)
		}
		if tc.want == 303 && rr.Header().Get("Location") != "/pair" {
			t.Error("missing pair redirect")
		}
		if tc.want == 401 && !strings.Contains(rr.Body.String(), "pairing_required") {
			t.Error("API must return pairing JSON, not HTML")
		}
	}
}
