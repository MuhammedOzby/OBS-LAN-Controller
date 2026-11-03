package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "configuration.json")
	if err := os.WriteFile(path, []byte(`{"obs":{"ip":"192.168.1.20","port":4455,"password":"test-only"},"server":{"listening_ip":"127.0.0.1","port":8090}}`), 0600); err != nil {
		t.Fatal(err)
	}
	config, err := ReadConfiguration(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.OBS.IP != "192.168.1.20" || config.OBS.Port != 4455 || config.OBS.Password != "test-only" {
		t.Error("OBS settings were not read correctly")
	}
	if config.Server.ListeningIP != "127.0.0.1" || config.Server.Port != 8090 {
		t.Error("server settings were not read correctly")
	}
}
