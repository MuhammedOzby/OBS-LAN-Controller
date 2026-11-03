package lib

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFavoritePersistenceAndFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prefs.json")
	s, err := loadFavorites(path)
	if err != nil {
		t.Fatal(err)
	}
	want := favoriteConfig{Collection: "Demo", Slots: [3]string{"Scene3", "", "Scene1"}}
	if err := s.save(want); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("mode: %v %v", info, err)
	}
	reloaded, err := loadFavorites(path)
	if err != nil || *reloaded.get() != want {
		t.Fatalf("reload: %v %v", reloaded, err)
	}
	s.path = filepath.Join(path, "impossible")
	if s.save(favoriteConfig{Collection: "Other"}) == nil || *s.get() != want {
		t.Fatal("failed write published")
	}
}

func TestFavoritesPairedLAN(t *testing.T) {
	_, obs, _ := deckFixture(t)
	store, _ := loadFavorites("")
	h, err := newHandlerWithFavorites(obs, nil, store)
	if err != nil {
		t.Fatal(err)
	}
	pair := newPairingHandler(h)
	code := pairingCode(t, pair)
	paired := pairingRequest(pair, "POST", "/pair", "192.0.2.3:5000", "http://localhost:8080", url.Values{"code": {code}}.Encode(), nil)
	if paired.Code != 303 || len(paired.Result().Cookies()) != 1 {
		t.Fatal("pair failed")
	}
	request := func(cookie *http.Cookie) int {
		r := httptest.NewRequest("POST", "http://localhost:8080/api/favorites", strings.NewReader(`{"collection":"Demo","slots":["Scene3","",""]}`))
		r.RemoteAddr = "192.0.2.3:5000"
		r.Header.Set("Origin", "http://localhost:8080")
		r.Header.Set("Content-Type", "application/json")
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		pair.ServeHTTP(w, r)
		return w.Code
	}
	if got := request(nil); got != 401 {
		t.Fatalf("unpaired %d", got)
	}
	if got := request(paired.Result().Cookies()[0]); got != 204 {
		t.Fatalf("paired %d", got)
	}
}

func TestFavoritesHTTPAndAuthorization(t *testing.T) {
	sim, obs, _ := deckFixture(t)
	store, _ := loadFavorites(filepath.Join(t.TempDir(), "prefs.json"))
	h, err := newHandlerWithFavorites(obs, nil, store)
	if err != nil {
		t.Fatal(err)
	}
	call := func(method, path, body, origin, content string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://controller"+path, strings.NewReader(body))
		r.Header.Set("Origin", origin)
		r.Header.Set("Content-Type", content)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	valid := `{"collection":"Demo","slots":["Scene3","","Scene1"]}`
	if w := call("POST", "/api/favorites", valid, "http://evil", "application/json"); w.Code != 403 {
		t.Fatal(w.Code)
	}
	if w := call("POST", "/api/favorites", valid, "http://controller", "text/plain"); w.Code != 415 {
		t.Fatal(w.Code)
	}
	if w := call("POST", "/api/favorites", `{"collection":"Demo","slots":["Scene3","Scene3",""]}`, "http://controller", "application/json"); w.Code != 400 {
		t.Fatal(w.Code)
	}
	if w := call("POST", "/api/favorites", valid, "http://controller", "application/json"); w.Code != 204 {
		t.Fatal(w.Code)
	}
	store.path = filepath.Join(store.path, "impossible")
	if w := call("POST", "/api/favorites", `{"collection":"Demo","slots":["Scene0","",""]}`, "http://controller", "application/json"); w.Code != 500 {
		t.Fatalf("write failure status %d", w.Code)
	}
	if got := store.get(); got.Slots[0] != "Scene3" {
		t.Fatal("failed write changed published selection")
	}
	if err := obs.SetScene(context.Background(), "Scene3"); err != nil {
		t.Fatal(err)
	}
	if err := obs.SetScene(context.Background(), "Scene0"); err == nil {
		t.Fatal("unselected scene authorized")
	}
	w := call("GET", "/api/favorites", "", "", "")
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Scene3") {
		t.Fatal(w.Code, w.Body.String())
	}
	sim.mu.Lock()
	sim.scenes = 3
	sim.mu.Unlock()
	if err := obs.SetScene(context.Background(), "Scene3"); err == nil {
		t.Fatal("deleted scene authorized")
	}
	w = call("GET", "/api/favorites", "", "", "")
	if !strings.Contains(w.Body.String(), "Scene3") {
		t.Fatal("lost choice remapped")
	}
	sim.mu.Lock()
	sim.collection = "Other"
	sim.mu.Unlock()
	if err := obs.SetScene(context.Background(), "Scene1"); err == nil {
		t.Fatal("wrong collection authorized")
	}
	if w := call("POST", "/api/favorites", valid, "http://controller", "application/json"); w.Code != 409 {
		t.Fatal(w.Code)
	}
}
