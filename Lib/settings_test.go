package lib

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSettingsSecurityPersistenceAndReconfigure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	var c Configuration
	c.OBS.IP, c.OBS.Port, c.OBS.Password = "127.0.0.1", 4455, "secret-value"
	c.Server.ListeningIP, c.Server.Port = "127.0.0.1", 8080
	obs := NewOBSClient(c)
	h := newSettingsHandler(http.NotFoundHandler(), obs, path, true)
	request := func(method, peer, body string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://localhost/api/admin/config", strings.NewReader(body))
		r.RemoteAddr = peer
		r.Header.Set("Origin", "http://localhost")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-Forwarded-For", "127.0.0.1")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := request("GET", "192.168.1.2:1000", ""); w.Code != 403 {
		t.Fatal(w.Code)
	}
	if w := request("GET", "127.0.0.1:1000", ""); w.Code != 200 || bytes.Contains(w.Body.Bytes(), []byte("secret-value")) {
		t.Fatal(w.Code, w.Body.String())
	}
	body := `{"host":"127.0.0.2","port":4456,"listeningIP":"0.0.0.0","listeningPort":8080}`
	obs.retryAfter = time.Now().Add(time.Hour)
	obs.cached.Connected = true
	if w := request("POST", "127.0.0.1:1000", body); w.Code != 204 {
		t.Fatal(w.Code, w.Body.String())
	}
	got, err := ReadConfiguration(path)
	if err != nil || got.OBS.Password != "secret-value" || got.OBS.IP != "127.0.0.2" {
		t.Fatal(got, err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0600 {
		t.Fatal(info.Mode())
	}
	if obs.cached.Connected || !obs.retryAfter.IsZero() {
		t.Fatal("stale connection state")
	}
	clear := strings.Replace(body, `"port":4456`, `"port":4456,"password":""`, 1)
	if w := request("POST", "127.0.0.1:1000", clear); w.Code != 204 {
		t.Fatal(w.Code)
	}
	if obs.config.OBS.Password != "" {
		t.Fatal("password not cleared")
	}
	if w := request("POST", "127.0.0.1:1000", body+`{}`); w.Code != 400 {
		t.Fatal(w.Code)
	}
	if w := request("POST", "127.0.0.1:1000", strings.Replace(body, `"port":4456`, `"port":4456,"extra":1`, 1)); w.Code != 400 {
		t.Fatal(w.Code)
	}
	bad := newSettingsHandler(http.NotFoundHandler(), obs, filepath.Join(t.TempDir(), "absent", "config"), false)
	r := httptest.NewRequest("POST", "http://localhost/api/admin/config", strings.NewReader(strings.Replace(body, "127.0.0.2", "127.0.0.3", 1)))
	r.RemoteAddr = "127.0.0.1:1"
	r.Header.Set("Origin", "http://localhost")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	bad.ServeHTTP(w, r)
	if w.Code != 500 || obs.config.OBS.IP != "127.0.0.2" {
		t.Fatal(w.Code, obs.config.OBS.IP)
	}
}

func TestSettingsLocalUIAndOrigin(t *testing.T) {
	h, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/deck", "/yayin-yonetimi"} {
		for _, local := range []bool{true, false} {
			r := httptest.NewRequest("GET", "http://localhost"+path, nil)
			if local {
				r.RemoteAddr = "127.0.0.1:1"
			} else {
				r.RemoteAddr = "192.168.1.5:1"
				r.Header.Set("X-Forwarded-For", "127.0.0.1")
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != 200 || strings.Contains(w.Body.String(), `id="settings-open"`) != local {
				t.Fatalf("%s local=%v status=%d", path, local, w.Code)
			}
		}
	}
	obs := NewOBSClient(Configuration{})
	s := newSettingsHandler(http.NotFoundHandler(), obs, filepath.Join(t.TempDir(), "config"), true)
	r := httptest.NewRequest("POST", "http://localhost/api/admin/config", strings.NewReader(`{}`))
	r.RemoteAddr = "127.0.0.1:1"
	r.Header.Set("Origin", "http://evil.example")
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
}

func TestRunRejectsMalformedAndLANBootstrap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte(`{broken`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"-config", path}); err == nil {
		t.Fatal("malformed configuration accepted")
	}
	if err := Run([]string{"-config", filepath.Join(dir, "missing"), "-addr", "0.0.0.0:8080"}); err == nil {
		t.Fatal("LAN bootstrap accepted")
	}
}
