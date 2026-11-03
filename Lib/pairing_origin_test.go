package lib

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPairingFormPreservesSameOriginOnPost(t *testing.T) {
	h := newPairingHandler(http.NotFoundHandler())
	for _, peer := range []string{"127.0.0.1:5000", "192.0.2.3:5000"} {
		rr := pairingRequest(h, "GET", "/pair", peer, "", "", nil)
		// Chromium native navigation forms send Origin: null under no-referrer.
		// Preserve same-origin provenance without leaking to other origins.
		if got := rr.Header().Get("Referrer-Policy"); got != "same-origin" {
			t.Fatalf("form policy = %q; want same-origin", got)
		}
	}
	for _, origin := range []string{"null", "http://attacker.example"} {
		if rr := pairingRequest(h, "POST", "/pair", "192.0.2.3:5000", origin, "code=000000000000", nil); rr.Code != 403 {
			t.Fatalf("untrusted Origin %q accepted: %d", origin, rr.Code)
		}
	}
}

func TestPairingCrossSiteScannerLanding(t *testing.T) {
	h := newPairingHandler(http.NotFoundHandler())
	code := pairingCode(t, h)
	request := func(method, peer, origin string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://localhost:8080/pair", strings.NewReader("code="+code))
		r.RemoteAddr = peer
		r.Header.Set("Sec-Fetch-Site", "cross-site")
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	for _, origin := range []string{"", "http://scanner.example"} {
		local := request("GET", "127.0.0.1:123", origin)
		if local.Code != 403 || h.code != code {
			t.Fatal("cross-site local issuance allowed")
		}
		remote := request("GET", "192.0.2.1:123", origin)
		if remote.Code != 200 || !strings.Contains(remote.Body.String(), `id="pair-form"`) || !strings.Contains(remote.Body.String(), "/pair-bootstrap.js") || strings.Contains(remote.Body.String(), code) || len(remote.Result().Cookies()) != 0 || h.code != code || len(h.sessions) != 0 {
			t.Fatal("remote scanner landing disclosed or changed state")
		}
	}
	if r := request("POST", "192.0.2.1:123", "http://localhost:8080"); r.Code != 403 || h.code != code || len(h.sessions) != 0 {
		t.Fatal("cross-site POST accepted")
	}
}
