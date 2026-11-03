package lib

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// All output mutations in these tests target this loopback simulator, never OBS.
type broadcastSimulator struct {
	mu                                            sync.Mutex
	stream, record, paused                        bool
	reject, disconnect                            string
	disconnectReadback                            bool
	mutated                                       bool
	recordUnavailable, emptyRecord, ignoreCommand bool
	delayed                                       int
	pending                                       string
	requests                                      []string
}

func broadcastFixture(t *testing.T) (*broadcastSimulator, *OBSClient, http.Handler) {
	t.Helper()
	sim := &broadcastSimulator{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		conn.WriteJSON(map[string]any{"op": 0, "d": map[string]any{"rpcVersion": 1}})
		var identify map[string]any
		if conn.ReadJSON(&identify) != nil {
			return
		}
		conn.WriteJSON(map[string]any{"op": 2, "d": map[string]any{"negotiatedRpcVersion": 1}})
		for {
			var req struct {
				D struct {
					Kind string `json:"requestType"`
					ID   string `json:"requestId"`
				} `json:"d"`
			}
			if conn.ReadJSON(&req) != nil {
				return
			}
			sim.mu.Lock()
			kind := req.D.Kind
			sim.requests = append(sim.requests, kind)
			if kind == sim.disconnect || (sim.disconnectReadback && sim.mutated && (kind == "GetStreamStatus" || kind == "GetRecordStatus")) {
				sim.mu.Unlock()
				return
			}
			result := kind != sim.reject && !(kind == "GetRecordStatus" && sim.recordUnavailable)
			data := map[string]any{}
			if result {
				switch kind {
				case "GetSceneList":
					data = map[string]any{"currentProgramSceneName": "Main", "scenes": []any{map[string]any{"sceneName": "Main"}}}
				case "GetStats":
					data["activeFps"] = 60
				case "GetStreamStatus", "GetRecordStatus":
					if sim.pending != "" {
						if sim.delayed > 0 {
							sim.delayed--
						} else {
							sim.apply(sim.pending)
							sim.pending = ""
						}
					}
					if kind == "GetStreamStatus" {
						data = map[string]any{"outputActive": sim.stream, "outputReconnecting": false, "outputDuration": 1200, "outputBytes": 12345}
					} else if !sim.emptyRecord {
						data = map[string]any{"outputActive": sim.record, "outputPaused": sim.paused, "outputDuration": 800, "outputBytes": 6789}
					}
				case "StartStream", "StopStream", "StartRecord", "StopRecord", "PauseRecord", "ResumeRecord":
					sim.mutated = true
					if !sim.ignoreCommand {
						if sim.delayed > 0 {
							sim.pending = kind
						} else {
							sim.apply(kind)
						}
					}
				}
			}
			sim.mu.Unlock()
			if conn.WriteJSON(map[string]any{"op": 7, "d": map[string]any{"requestId": req.D.ID, "requestStatus": map[string]any{"result": result, "code": 500}, "responseData": data}}) != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	host, port, _ := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	n, _ := strconv.Atoi(port)
	cfg := Configuration{}
	cfg.OBS.IP = host
	cfg.OBS.Port = n
	cfg.OBS.Password = "fixture-private-secret"
	client := NewOBSClient(cfg)
	t.Cleanup(client.Close)
	handler, err := newHandler(client)
	if err != nil {
		t.Fatal(err)
	}
	return sim, client, handler
}
func (s *broadcastSimulator) apply(kind string) {
	switch kind {
	case "StartStream":
		s.stream = true
	case "StopStream":
		s.stream = false
	case "StartRecord":
		s.record = true
		s.paused = false
	case "StopRecord":
		s.record = false
		s.paused = false
	case "PauseRecord":
		s.paused = true
	case "ResumeRecord":
		s.paused = false
	}
}
func broadcastPost(h http.Handler, body, origin string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "http://controller.test/api/broadcast", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
func broadcastJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("JSON %d %q: %v", rr.Code, rr.Body, err)
	}
	return got
}
func TestBroadcastSnapshotRecording(t *testing.T) {
	sim, client, _ := broadcastFixture(t)
	sim.record = true
	sim.paused = true
	raw, _ := json.Marshal(client.Snapshot(context.Background(), false))
	var got map[string]any
	json.Unmarshal(raw, &got)
	for key, want := range map[string]any{"connected": true, "recordAvailable": true, "recordActive": true, "recordPaused": true, "recordDuration": float64(800), "recordBytes": float64(6789), "streamBytes": float64(12345)} {
		if got[key] != want {
			t.Errorf("%s = %v, want %v", key, got[key], want)
		}
	}
}
func TestBroadcastStartReadback(t *testing.T) {
	sim, _, handler := broadcastFixture(t)
	sim.delayed = 6
	rr := broadcastPost(handler, `{"action":"start-stream","confirmed":true}`, "http://controller.test")
	if rr.Code != 200 {
		t.Fatalf("start = %d %s", rr.Code, rr.Body)
	}
	got := broadcastJSON(t, rr)
	if got["streamActive"] != true || got["recordAvailable"] != true || got["streamBytes"] != float64(12345) || got["recordBytes"] != float64(6789) {
		t.Fatalf("fresh projection = %v", got)
	}
	if strings.Contains(rr.Body.String(), "fixture-private-secret") {
		t.Fatal("credentials leaked")
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	calls := strings.Join(sim.requests, ",")
	if !strings.Contains(calls, "GetStreamStatus,StartStream,GetStreamStatus,GetStreamStatus") {
		t.Fatalf("missing preflight/readback: %s", calls)
	}
}

func TestBroadcastActionsAndNoOps(t *testing.T) {
	for _, tc := range []struct {
		action, command, field       string
		stream, record, paused, want bool
	}{
		{"start-stream", "StartStream", "streamActive", false, false, false, true},
		{"stop-stream", "StopStream", "streamActive", true, false, false, false},
		{"start-record", "StartRecord", "recordActive", false, false, false, true},
		{"stop-record", "StopRecord", "recordActive", false, true, true, false},
		{"pause-record", "PauseRecord", "recordPaused", false, true, false, true},
		{"resume-record", "ResumeRecord", "recordPaused", false, true, true, false},
	} {
		t.Run(tc.action, func(t *testing.T) {
			sim, client, h := broadcastFixture(t)
			sim.stream = tc.stream
			sim.record = tc.record
			sim.paused = tc.paused
			client.Snapshot(context.Background(), false)
			body := `{"action":"` + tc.action + `","confirmed":true}`
			rr := broadcastPost(h, body, "http://controller.test")
			if rr.Code != 200 {
				t.Fatalf("command %d %s", rr.Code, rr.Body)
			}
			if got := broadcastJSON(t, rr); got[tc.field] != tc.want {
				t.Fatalf("projection %v", got)
			}
			rr = broadcastPost(h, body, "http://controller.test")
			if rr.Code != 409 {
				t.Fatalf("no-op %d", rr.Code)
			}
			sim.mu.Lock()
			defer sim.mu.Unlock()
			n := 0
			for _, kind := range sim.requests {
				if kind == tc.command {
					n++
				}
				if strings.Contains(kind, "Toggle") {
					t.Fatal("unsafe toggle")
				}
			}
			if n != 1 {
				t.Fatalf("command count %d", n)
			}
		})
	}
}
func TestBroadcastConcurrentDoubleSubmit(t *testing.T) {
	sim, _, h := broadcastFixture(t)
	codes := make(chan int, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes <- broadcastPost(h, `{"action":"start-record","confirmed":true}`, "http://controller.test").Code
		}()
	}
	wg.Wait()
	close(codes)
	counts := map[int]int{}
	for code := range codes {
		counts[code]++
	}
	if counts[200] != 1 || counts[409] != 1 {
		t.Fatalf("codes %v", counts)
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	n := 0
	for _, kind := range sim.requests {
		if kind == "StartRecord" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("mutations %d", n)
	}
}
func TestBroadcastFreshPreflight(t *testing.T) {
	sim, client, h := broadcastFixture(t)
	client.Snapshot(context.Background(), false)
	sim.mu.Lock()
	sim.stream = true
	sim.mu.Unlock()
	rr := broadcastPost(h, `{"action":"start-stream","confirmed":true}`, "http://controller.test")
	if rr.Code != 409 {
		t.Fatalf("stale display status used: %d", rr.Code)
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	for _, kind := range sim.requests {
		if kind == "StartStream" {
			t.Fatal("stale command sent")
		}
	}
}
func TestBroadcastUnavailableRecording(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(strconv.FormatBool(empty), func(t *testing.T) {
			sim, client, h := broadcastFixture(t)
			sim.emptyRecord = empty
			sim.recordUnavailable = !empty
			raw, _ := json.Marshal(client.Snapshot(context.Background(), false))
			var got map[string]any
			json.Unmarshal(raw, &got)
			if got["connected"] != true || got["recordAvailable"] != false {
				t.Fatalf("unsupported record disconnected: %s", raw)
			}
			rr := broadcastPost(h, `{"action":"start-record","confirmed":true}`, "http://controller.test")
			if rr.Code != 409 {
				t.Fatalf("unsupported %d", rr.Code)
			}
			rr = broadcastPost(h, `{"action":"start-stream","confirmed":true}`, "http://controller.test")
			if rr.Code != 200 {
				t.Fatalf("stream with unsupported recording %d", rr.Code)
			}
		})
	}
}
func TestBroadcastRejectedAndUnverified(t *testing.T) {
	for _, tc := range []struct {
		name   string
		reject string
		ignore bool
		want   int
	}{
		{"rejection", "StartStream", false, 409}, {"pause unsupported", "PauseRecord", false, 409}, {"ack without state", "", true, 502},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sim, client, h := broadcastFixture(t)
			sim.reject = tc.reject
			sim.ignoreCommand = tc.ignore
			action := "start-stream"
			if tc.reject == "PauseRecord" {
				sim.record = true
				action = "pause-record"
			}
			start := time.Now()
			rr := broadcastPost(h, `{"action":"`+action+`","confirmed":true}`, "http://controller.test")
			if rr.Code != tc.want {
				t.Fatalf("got %d want %d", rr.Code, tc.want)
			}
			if time.Since(start) > 4*time.Second {
				t.Fatal("verification not bounded")
			}
			if !client.Snapshot(context.Background(), false).Connected {
				t.Fatal("rejection disconnected client")
			}
			sim.mu.Lock()
			defer sim.mu.Unlock()
			n := 0
			for _, kind := range sim.requests {
				if kind == "StartStream" || kind == "PauseRecord" {
					n++
				}
			}
			if n != 1 {
				t.Fatalf("command attempts %d", n)
			}
		})
	}
}
func TestBroadcastDisconnectNeverReplays(t *testing.T) {
	sim, client, h := broadcastFixture(t)
	sim.disconnect = "StartStream"
	rr := broadcastPost(h, `{"action":"start-stream","confirmed":true}`, "http://controller.test")
	if rr.Code != 502 {
		t.Fatalf("disconnect %d", rr.Code)
	}
	// Allow only the ordinary read-only poll to reconnect; never resend the action.
	sim.mu.Lock()
	sim.disconnect = ""
	sim.mu.Unlock()
	client.mu.Lock()
	client.retryAfter = time.Time{}
	client.mu.Unlock()
	if !client.Snapshot(context.Background(), false).Connected {
		t.Fatal("reconnect failed")
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	n := 0
	for _, kind := range sim.requests {
		if kind == "StartStream" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("replayed %d", n)
	}
}
func TestBroadcastReadbackDisconnectNeverReplays(t *testing.T) {
	sim, client, h := broadcastFixture(t)
	sim.disconnectReadback = true
	rr := broadcastPost(h, `{"action":"start-record","confirmed":true}`, "http://controller.test")
	if rr.Code != 502 {
		t.Fatalf("readback disconnect %d", rr.Code)
	}
	sim.mu.Lock()
	sim.disconnectReadback = false
	sim.mu.Unlock()
	client.mu.Lock()
	client.retryAfter = time.Time{}
	client.mu.Unlock()
	status := client.Snapshot(context.Background(), false)
	if !status.Connected || !status.RecordActive {
		t.Fatalf("reconnect did not recover actual state: %+v", status)
	}
	rr = broadcastPost(h, `{"action":"start-record","confirmed":true}`, "http://controller.test")
	if rr.Code != 409 {
		t.Fatalf("post-reconnect duplicate %d", rr.Code)
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	n := 0
	for _, kind := range sim.requests {
		if kind == "StartRecord" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("replayed command %d", n)
	}
}

func TestBroadcastValidationAndOrigin(t *testing.T) {
	sim, _, h := broadcastFixture(t)
	for _, origin := range []string{"", "null", "http://evil.test", "https://controller.test", "http://u@controller.test", "http://controller.test/", "http://controller.test?x=1", "http://controller.test#x"} {
		if rr := broadcastPost(h, `{"action":"start-stream","confirmed":true}`, origin); rr.Code != 403 {
			t.Errorf("origin %q = %d", origin, rr.Code)
		}
	}
	for _, body := range []string{`{}`, `{"action":"start-stream"}`, `{"action":"start-stream","confirmed":false}`, `{"action":"start-stream","confirmed":"true"}`, `{"action":"ToggleStream","confirmed":true}`, `{"action":"SetStreamServiceSettings","confirmed":true}`, `{"action":"start-stream","confirmed":true,"extra":1}`, `{"action":"start-stream","confirmed":true} {}`, `{"action":"start-stream","confirmed":true} garbage`, strings.Repeat("x", 4097)} {
		if rr := broadcastPost(h, body, "http://controller.test"); rr.Code != 400 {
			t.Errorf("invalid body = %d: %.100s", rr.Code, body)
		}
	}
	sim.mu.Lock()
	defer sim.mu.Unlock()
	if len(sim.requests) != 0 {
		t.Fatalf("invalid request reached OBS: %v", sim.requests)
	}
}
func TestBroadcastOffline(t *testing.T) {
	h, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	if rr := broadcastPost(h, `{"action":"start-stream","confirmed":true}`, "http://controller.test"); rr.Code != 503 {
		t.Fatalf("nil client %d", rr.Code)
	}
	sim, _, h := broadcastFixture(t)
	sim.disconnect = "GetStreamStatus"
	if rr := broadcastPost(h, `{"action":"start-stream","confirmed":true}`, "http://controller.test"); rr.Code != 503 {
		t.Fatalf("preflight offline %d", rr.Code)
	}
}
