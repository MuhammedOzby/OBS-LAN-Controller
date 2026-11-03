package lib

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestReadConfigurationRejectsInvalidSettings(t *testing.T) {
	for _, input := range []string{
		`{`, `null`, `{}`,
		`{"obs":{"ip":"invalid","port":4455},"server":{"listening_ip":"127.0.0.1","port":8080}}`,
		`{"obs":{"ip":"127.0.0.1","port":0},"server":{"listening_ip":"127.0.0.1","port":8080}}`,
		`{"obs":{"ip":"127.0.0.1","port":65536},"server":{"listening_ip":"127.0.0.1","port":8080}}`,
		`{"obs":{"ip":"127.0.0.1","port":4455},"server":{"listening_ip":"localhost","port":8080}}`,
		`{"obs":{"ip":"127.0.0.1","port":4455},"server":{"listening_ip":"127.0.0.1","port":-1}}`,
		`{"obs":{"ip":"127.0.0.1","port":4455},"server":{"listening_ip":"127.0.0.1","port":65536}}`,
		`{"obs":{"ip":"127.0.0.1","port":4455,"password":{"test-only":true}},"server":{"listening_ip":"127.0.0.1","port":8080}}`,
	} {
		path := filepath.Join(t.TempDir(), "configuration.json")
		if err := os.WriteFile(path, []byte(input), 0600); err != nil {
			t.Fatal(err)
		}
		config, err := ReadConfiguration(path)
		if err == nil {
			t.Error("invalid configuration accepted")
		} else if strings.Contains(err.Error(), "test-only") {
			t.Error("error exposes configuration contents")
		}
		if config != (Configuration{}) {
			t.Error("invalid configuration returned partially loaded settings")
		}
	}
	if _, err := ReadConfiguration(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Error("missing file accepted")
	}
}

func TestRunUsesConfiguredListeningAddress(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	path := filepath.Join(t.TempDir(), "configuration.json")
	input := `{"obs":{"ip":"127.0.0.1","port":4455,"password":"test-only"},"server":{"listening_ip":"127.0.0.1","port":` + strconv.Itoa(port) + `}}`
	if err := os.WriteFile(path, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	err = Run([]string{"serve", "-config", path})
	if err == nil || !strings.Contains(err.Error(), listener.Addr().String()) {
		t.Fatalf("expected bind failure at configured address, got %v", err)
	}
}
