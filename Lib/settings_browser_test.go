package lib

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSettingsProductionChainBrowser(t *testing.T) {
	chrome := os.Getenv("CHROME_BIN")
	if chrome == "" {
		t.Skip("CHROME_BIN required")
	}
	var c Configuration
	c.OBS.IP, c.OBS.Port = "127.0.0.1", 4455
	c.Server.ListeningIP, c.Server.Port = "127.0.0.1", 8080
	obs := NewOBSClient(c)
	defer obs.Close()
	router, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	fixture := `<script>
 window.addEventListener('load', async()=>{try {
 const q=id=>document.getElementById('settings-'+id), delay=()=>new Promise(r=>setTimeout(r,50));
 const wait=async f=>{for(let i=0;i<100;i++){if(f())return;await delay()}throw Error('timeout')};
 const state=async()=>await (await fetch('/test/state')).text();
 const save=async()=>{q('form').requestSubmit(); await wait(()=>!q('form').querySelector('[type=submit]').disabled);};
 await wait(()=>q('dialog').open); if(q('password').value)throw Error('initial secret');
 q('password').value='first-test-secret'; await save(); if(await state()!=='first')throw Error('bootstrap');
 q('close').click(); q('open').click(); await wait(()=>q('dialog').open); if(q('password').value)throw Error('secret exposed');
 q('host').value='invalid'; await save(); if(!q('message').textContent.includes('kaydedilemedi') || q('host').disabled || q('host').value!=='invalid')throw Error('failed save');
 q('host').value='127.0.0.1'; await save(); if(await state()!=='first')throw Error('preserve');
 q('password').value='second-test-secret'; await save(); if(await state()!=='second'||q('password').value)throw Error('replace');
 q('close').click(); q('open').click(); await wait(()=>q('dialog').open); if(q('password').value)throw Error('replacement exposed');
 q('clear').checked=true; await save(); if(await state()!=='empty')throw Error('clear');
 document.body.dataset.settingsResult='PASS';
 }catch(e){document.body.dataset.settingsResult='FAIL '+e.message}});
 </script>`
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test/state" {
			obs.mu.Lock()
			defer obs.mu.Unlock()
			switch obs.config.OBS.Password {
			case "first-test-secret":
				fmt.Fprint(w, "first")
			case "second-test-secret":
				fmt.Fprint(w, "second")
			case "":
				fmt.Fprint(w, "empty")
			default:
				fmt.Fprint(w, "other")
			}
			return
		}
		if r.URL.Path == "/api/status" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"connected":false,"scenes":[]}`)
			return
		}
		if r.URL.Path == "/" {
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, r)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, strings.Replace(rec.Body.String(), "<head>", "<head>"+fixture, 1))
			return
		}
		router.ServeHTTP(w, r)
	})
	path := filepath.Join(t.TempDir(), "config.json")
	pairing := newPairingHandler(newSettingsHandler(base, obs, path, true))
	server, err := newHTTPServer("127.0.0.1:0", pairing)
	if err != nil {
		t.Fatal(err)
	}
	// Exercise the same guard, pairing and settings chain with a paired non-loopback peer.
	token := strings.Repeat("a", 64)
	pairing.sessions[token] = time.Now().Add(time.Hour)
	for _, tc := range []struct {
		path, host, origin, peer string
		want                     int
	}{{"/api/admin/config", "localhost", "http://localhost", "192.0.2.1:123", 403}, {"/api/admin/config", "evil.example", "http://evil.example", "127.0.0.1:1", 403}, {"/api/admin/config", "localhost", "http://evil.example", "127.0.0.1:1", 403}, {"/", "localhost", "", "192.0.2.1:1", 200}} {
		r := httptest.NewRequest("GET", "http://"+tc.host+tc.path, nil)
		r.RemoteAddr = tc.peer
		r.Header.Set("X-Forwarded-For", "127.0.0.1")
		r.Header.Set("X-Real-IP", "127.0.0.1")
		if tc.origin != "" {
			r.Header.Set("Origin", tc.origin)
		}
		r.AddCookie(&http.Cookie{Name: pairingCookieName, Value: token})
		w := httptest.NewRecorder()
		server.Handler.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%s: %d want %d", tc.path, w.Code, tc.want)
		}
		if tc.path == "/" && (strings.Contains(w.Body.String(), `id="settings-open"`) || strings.Contains(w.Body.String(), "settings.js")) {
			t.Fatal("remote admin UI exposed")
		}
	}
	for _, tc := range []struct{ peer, origin string }{{"192.0.2.1:123", "http://localhost"}, {"127.0.0.1:1", "http://evil.example"}} {
		r := httptest.NewRequest("POST", "http://localhost/api/admin/config", strings.NewReader(`{}`))
		r.RemoteAddr = tc.peer
		r.Header.Set("Origin", tc.origin)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("X-Forwarded-For", "127.0.0.1")
		r.AddCookie(&http.Cookie{Name: pairingCookieName, Value: token})
		w := httptest.NewRecorder()
		server.Handler.ServeHTTP(w, r)
		if w.Code != http.StatusForbidden {
			t.Fatalf("admin POST peer %s: %d", tc.peer, w.Code)
		}
	}
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	out, err := exec.Command(chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--user-data-dir="+t.TempDir(), "--virtual-time-budget=10000", "--dump-dom", ts.URL).CombinedOutput()
	if err != nil || !bytes.Contains(out, []byte(`data-settings-result="PASS"`)) {
		t.Fatalf("browser: %v\n%s", err, out)
	}
	got, err := ReadConfiguration(path)
	if err != nil || got.OBS.Password != "" {
		t.Fatalf("final config invalid: %v", err)
	}
}
