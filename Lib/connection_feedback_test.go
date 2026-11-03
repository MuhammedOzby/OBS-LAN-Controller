package lib

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"
)

func TestConnectionFeedbackDoesNotExposeConfiguration(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	var config Configuration
	config.OBS.IP, config.OBS.Port, config.OBS.Password = "127.0.0.1", port, "test-secret-never-expose"
	client := NewOBSClient(config)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	status := client.Snapshot(ctx, false)
	if status.Connected || status.ConnectionError == "" {
		t.Fatalf("unexpected status: %+v", status)
	}
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), config.OBS.Password) {
		t.Fatal("secret exposed in status")
	}
	status = client.Snapshot(ctx, false)
	if !strings.Contains(status.ConnectionError, "yeniden") {
		t.Fatalf("retry feedback: %q", status.ConnectionError)
	}
}
