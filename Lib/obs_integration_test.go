package lib

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestOBSAuthenticatedSnapshotAndScene(t *testing.T) {
	var mu sync.Mutex
	current := "Scene A"
	commands := []string{}
	authOK := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		conn.WriteJSON(map[string]any{"op": 0, "d": map[string]any{"rpcVersion": 1, "authentication": map[string]any{"salt": "test-salt", "challenge": "test-challenge"}}})
		var identify struct {
			D struct {
				Authentication string `json:"authentication"`
				Events         int    `json:"eventSubscriptions"`
			} `json:"d"`
		}
		if conn.ReadJSON(&identify) != nil {
			return
		}
		secret := sha256.Sum256([]byte("test-passwordtest-salt"))
		expected := sha256.Sum256([]byte(base64.StdEncoding.EncodeToString(secret[:]) + "test-challenge"))
		mu.Lock()
		authOK = identify.D.Authentication == base64.StdEncoding.EncodeToString(expected[:]) && identify.D.Events == 0
		mu.Unlock()
		conn.WriteJSON(map[string]any{"op": 2, "d": map[string]any{"negotiatedRpcVersion": 1}})
		for {
			var request struct {
				D struct {
					Kind string         `json:"requestType"`
					ID   string         `json:"requestId"`
					Args map[string]any `json:"requestData"`
				} `json:"d"`
			}
			if conn.ReadJSON(&request) != nil {
				return
			}
			mu.Lock()
			commands = append(commands, request.D.Kind)
			data := map[string]any{}
			switch request.D.Kind {
			case "GetSceneList":
				data = map[string]any{"currentProgramSceneName": current, "scenes": []any{map[string]any{"sceneName": "Scene A"}, map[string]any{"sceneName": "Scene B"}, map[string]any{"sceneName": "Scene C"}, map[string]any{"sceneName": "Scene D"}}}
			case "GetStats":
				data = map[string]any{"cpuUsage": 12.5, "activeFps": 30, "renderSkippedFrames": 2}
			case "GetStreamStatus":
				data = map[string]any{"outputActive": true, "outputBytes": 1000, "outputDuration": 20000, "outputTotalFrames": 600}
			case "GetSourceScreenshot":
				data = map[string]any{"imageData": "data:image/jpeg;base64,dGVzdA=="}
			case "SetCurrentProgramScene":
				current = request.D.Args["sceneName"].(string)
			case "GetCurrentProgramScene":
				data = map[string]any{"currentProgramSceneName": current}
			}
			mu.Unlock()
			if conn.WriteJSON(map[string]any{"op": 7, "d": map[string]any{"requestId": request.D.ID, "requestStatus": map[string]any{"result": true, "code": 100}, "responseData": data}}) != nil {
				return
			}
		}
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(u.Port())
	var config Configuration
	config.OBS.IP = u.Hostname()
	config.OBS.Port = port
	config.OBS.Password = "test-password"
	client := NewOBSClient(config)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status := client.Snapshot(ctx, true)
	if !status.Connected || !status.StreamActive || status.CPU != 12.5 || status.CurrentScene != "Scene A" || !status.PreviewAvailable || status.BitrateKbps != nil {
		t.Fatalf("snapshot %+v", status)
	}
	if got := client.Snapshot(ctx, false); got.Preview != "" || got.PreviewAvailable {
		t.Fatal("non preview response leaked screenshot")
	}
	if err := client.SetScene(ctx, "Scene D"); err == nil {
		t.Fatal("fourth scene allowed")
	}
	handler, err := newHandler(client)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "http://controller/api/scene", strings.NewReader(`{"sceneName":"Scene B"}`))
	req.Header.Set("Origin", "http://controller")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 204 {
		t.Fatalf("scene status %d", rr.Code)
	}
	got := client.Snapshot(ctx, true)
	if got.CurrentScene != "Scene B" {
		t.Fatalf("scene was not read back: %s", got.CurrentScene)
	}
	raw, _ := json.Marshal(got)
	if strings.Contains(string(raw), "test-password") {
		t.Fatal("secret exposed")
	}
	mu.Lock()
	defer mu.Unlock()
	if !authOK {
		t.Fatal("authentication mismatch")
	}
	found := false
	for _, kind := range commands {
		if kind == "GetCurrentProgramScene" {
			found = true
		}
		if strings.Contains(kind, "StartStream") || strings.Contains(kind, "StopStream") {
			t.Fatal("critical command")
		}
	}
	if !found {
		t.Fatal("scene readback missing")
	}
}
