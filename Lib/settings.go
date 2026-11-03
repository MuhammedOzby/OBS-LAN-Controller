package lib

import (
	"encoding/json"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// settingsHandler is local administration; authentication through a LAN session
// does not grant access. OBS's mutex also serializes persistence and reconfiguration.
type settingsHandler struct {
	next      http.Handler
	obs       *OBSClient
	path      string
	bootstrap bool
}

func newSettingsHandler(next http.Handler, obs *OBSClient, path string, bootstrap bool) http.Handler {
	return &settingsHandler{next: next, obs: obs, path: path, bootstrap: bootstrap}
}

func (s *settingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/admin/config" {
		s.next.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	if !pairingLoopback(r) {
		http.Error(w, "Yerel erişim gerekli", http.StatusForbidden)
		return
	}
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" || (r.Method == http.MethodPost || len(r.Header.Values("Origin")) != 0) && !pairingSameOrigin(r) {
		http.Error(w, "Geçersiz Origin", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodGet {
		s.obs.mu.Lock()
		defer s.obs.mu.Unlock()
		c := s.obs.config
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"host": c.OBS.IP, "port": c.OBS.Port, "passwordConfigured": c.OBS.Password != "", "listeningIP": c.Server.ListeningIP, "listeningPort": c.Server.Port, "firstSetup": s.bootstrap})
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Yöntem desteklenmiyor", 405)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		http.Error(w, "JSON gerekli", 415)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var v struct {
		Host          string  `json:"host"`
		Port          int     `json:"port"`
		Password      *string `json:"password"`
		ListeningIP   string  `json:"listeningIP"`
		ListeningPort int     `json:"listeningPort"`
	}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if d.Decode(&v) != nil || d.Decode(&struct{}{}) != io.EOF || net.ParseIP(v.Host) == nil || net.ParseIP(v.ListeningIP) == nil || v.Port < 1 || v.Port > 65535 || v.ListeningPort < 1 || v.ListeningPort > 65535 || v.Password != nil && len(*v.Password) > 1024 {
		http.Error(w, "Geçersiz ayarlar", 400)
		return
	}
	s.obs.mu.Lock()
	defer s.obs.mu.Unlock()
	c := s.obs.config
	c.OBS.IP, c.OBS.Port = v.Host, v.Port
	if v.Password != nil {
		c.OBS.Password = *v.Password
	}
	c.Server.ListeningIP, c.Server.Port = v.ListeningIP, v.ListeningPort
	if err := writeSettings(s.path, c, s.bootstrap); err != nil {
		http.Error(w, "Ayarlar kaydedilemedi", 500)
		return
	}
	s.obs.disconnect()
	s.obs.retryAfter = time.Time{}
	s.obs.previousBytes = 0
	s.obs.config = c
	s.bootstrap = false
	w.WriteHeader(http.StatusNoContent)
}

func writeSettings(path string, c Configuration, first bool) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".config-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if err = f.Chmod(0600); err == nil {
		_, err = f.Write(data)
	}
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if first {
		// Link is atomic and refuses to overwrite a configuration created meanwhile.
		return os.Link(name, path)
	}
	return os.Rename(name, path)
}
