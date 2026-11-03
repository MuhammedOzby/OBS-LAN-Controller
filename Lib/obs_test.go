package lib

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOfflineStatusAndSceneProtection(t *testing.T) {
	client := NewOBSClient(Configuration{})
	handler, err := newHandler(client)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/api/status?preview=0", nil))
	var status map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &status); err != nil {
		t.Fatal(err)
	}
	if rr.Code != 200 || status["connected"] != false {
		t.Fatalf("offline response %d %s", rr.Code, rr.Body)
	}
	req := httptest.NewRequest("POST", "/api/scene", nil)
	req.Header.Set("Origin", "http://evil.example")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("cross origin = %d", rr.Code)
	}
}

func TestBitrateCounterResets(t *testing.T) {
	now := time.Now()
	for _, tc := range []struct {
		active            bool
		previous, current float64
		seconds           float64
		valid             bool
		want              float64
	}{
		{true, 1000, 3000, 2, true, 8}, {false, 1000, 3000, 2, false, 0}, {true, 3000, 1000, 2, false, 0}, {true, 1000, 3000, 0, false, 0},
	} {
		got := bitrate(tc.active, tc.previous, tc.current, now, now.Add(time.Duration(tc.seconds)*time.Second))
		if (got != nil) != tc.valid || got != nil && *got != tc.want {
			t.Fatalf("bitrate = %v for %+v", got, tc)
		}
	}
}

func TestCancelledSnapshot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := NewOBSClient(Configuration{})
	if got := client.Snapshot(ctx, false); got.Connected {
		t.Fatal("cancelled client connected")
	}
}
