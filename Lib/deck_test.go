package lib

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

// Deck mutations only reach this real loopback websocket fixture, never live OBS.
type deckSimulator struct {
	mu                                   sync.Mutex
	current                              string
	enabled, muted                       bool
	media                                string
	cursor                               float64
	collection, sourceName               string
	scenes, items                        int
	reject, disconnect                   string
	ignore, disconnectAfter, revertFinal bool
	mutated                              bool
	delayed                              int
	pendingKind                          string
	pendingArgs                          map[string]any
	requests                             []string
	arguments                            []map[string]any
}

func deckFixture(t *testing.T) (*deckSimulator, *OBSClient, http.Handler) {
	t.Helper()
	s := &deckSimulator{current: "Scene0", media: "OBS_MEDIA_STATE_PAUSED", cursor: 10000, collection: "Demo", sourceName: "Camera", scenes: 4, items: 1}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		if c.WriteJSON(map[string]any{"op": 0, "d": map[string]any{"rpcVersion": 1}}) != nil {
			return
		}
		var identify any
		if c.ReadJSON(&identify) != nil {
			return
		}
		if c.WriteJSON(map[string]any{"op": 2, "d": map[string]any{"negotiatedRpcVersion": 1}}) != nil {
			return
		}
		for {
			var req struct {
				D struct {
					Kind string         `json:"requestType"`
					ID   string         `json:"requestId"`
					Args map[string]any `json:"requestData"`
				} `json:"d"`
			}
			if c.ReadJSON(&req) != nil {
				return
			}
			s.mu.Lock()
			k := req.D.Kind
			a := req.D.Args
			s.requests = append(s.requests, k)
			s.arguments = append(s.arguments, a)
			if k == s.disconnect || (s.disconnectAfter && s.mutated) {
				s.mu.Unlock()
				return
			}
			ok := k != s.reject
			data := map[string]any{}
			if ok {
				if s.pendingKind != "" && strings.HasPrefix(k, "Get") {
					if s.delayed > 0 {
						s.delayed--
					} else {
						s.apply(s.pendingKind, s.pendingArgs)
						s.pendingKind = ""
					}
				}
				switch k {
				case "GetSceneCollectionList":
					data["currentSceneCollectionName"] = s.collection
				case "GetStats":
					data["activeFps"] = 60
				case "GetStreamStatus", "GetRecordStatus":
					data["outputActive"] = false
				case "GetSceneList":
					if s.revertFinal && s.mutated {
						s.current = "Scene0"
					}
					scenes := []any{}
					for i := 0; i < s.scenes; i++ {
						scenes = append(scenes, map[string]any{"sceneName": fmt.Sprintf("Scene%d", i)})
					}
					data = map[string]any{"currentProgramSceneName": s.current, "scenes": scenes}
				case "GetCurrentProgramScene":
					data["currentProgramSceneName"] = s.current
				case "GetSceneItemList":
					items := []any{}
					for i := 0; i < s.items; i++ {
						items = append(items, map[string]any{"sceneItemId": i, "sourceName": s.sourceName, "sceneItemEnabled": s.enabled})
					}
					data["sceneItems"] = items
				case "GetSceneItemEnabled":
					data["sceneItemEnabled"] = s.enabled
				case "GetInputList":
					data["inputs"] = []any{map[string]any{"inputName": "Camera", "inputKind": "ffmpeg_source"}}
				case "GetInputMute":
					data["inputMuted"] = s.muted
				case "GetMediaInputStatus":
					data["mediaState"] = s.media
					data["mediaCursor"] = s.cursor
				case "SetCurrentProgramScene", "SetSceneItemEnabled", "SetInputMute", "TriggerMediaInputAction":
					s.mutated = true
					if !s.ignore {
						if s.delayed > 0 {
							s.pendingKind = k
							s.pendingArgs = a
						} else {
							s.apply(k, a)
						}
					}
				default:
					ok = false
				}
			}
			s.mu.Unlock()
			if c.WriteJSON(map[string]any{"op": 7, "d": map[string]any{"requestId": req.D.ID, "requestStatus": map[string]any{"result": ok, "code": 600, "comment": "private-server-detail"}, "responseData": data}}) != nil {
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
	client := NewOBSClient(cfg)
	t.Cleanup(client.Close)
	h, err := newHandler(client)
	if err != nil {
		t.Fatal(err)
	}
	return s, client, h
}
func (s *deckSimulator) apply(k string, a map[string]any) {
	switch k {
	case "SetCurrentProgramScene":
		s.current, _ = a["sceneName"].(string)
	case "SetSceneItemEnabled":
		s.enabled, _ = a["sceneItemEnabled"].(bool)
	case "SetInputMute":
		s.muted, _ = a["inputMuted"].(bool)
	case "TriggerMediaInputAction":
		switch a["mediaAction"] {
		case "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART":
			s.cursor = 0
			s.media = "OBS_MEDIA_STATE_PLAYING"
		case "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PLAY":
			s.media = "OBS_MEDIA_STATE_PLAYING"
		case "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PAUSE":
			s.media = "OBS_MEDIA_STATE_PAUSED"
		case "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_STOP":
			s.media = "OBS_MEDIA_STATE_STOPPED"
		}
	}
}
func deckHTTP(h http.Handler, method, body, origin string) *httptest.ResponseRecorder {
	path := "/api/deck"
	if method == "POST" {
		path += "/action"
	}
	req := httptest.NewRequest(method, "http://controller.test"+path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}
func TestDeckDiscovery(t *testing.T) {
	s, _, h := deckFixture(t)
	rr := deckHTTP(h, "GET", "", "")
	if rr.Code != 200 {
		t.Fatalf("GET = %d %s", rr.Code, rr.Body)
	}
	var got struct {
		Connected bool             `json:"connected"`
		Current   string           `json:"currentScene"`
		Scenes    []string         `json:"scenes"`
		Items     []map[string]any `json:"items"`
		Inputs    []map[string]any `json:"inputs"`
	}
	if json.Unmarshal(rr.Body.Bytes(), &got) != nil {
		t.Fatal(rr.Body.String())
	}
	if !got.Connected || got.Current != "Scene0" || len(got.Scenes) != 4 || len(got.Items) != 4 || len(got.Inputs) != 1 {
		t.Fatalf("discovery %+v", got)
	}
	if got.Items[3]["sceneName"] != "Scene3" || got.Items[3]["sceneItemId"] != float64(0) || got.Items[3]["enabled"] != false || got.Inputs[0]["muted"] != false || got.Inputs[0]["mediaState"] != "OBS_MEDIA_STATE_PAUSED" {
		t.Fatalf("schema %+v", got)
	}
	if rr.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("cacheable deck")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range s.requests {
		if !strings.HasPrefix(k, "Get") {
			t.Fatalf("GET mutation %s", k)
		}
	}
}

func TestDeckValidation(t *testing.T) {
	s, _, h := deckFixture(t)
	for _, body := range []string{
		`{}`, `null`, `[]`, `{"action":"StartStream"}`, `{"action":"start-stream"}`, `{"action":"macro","inputName":"Camera"}`,
		`{"action":"scene","sceneName":"Scene1","extra":true}`, `{"action":"scene","sceneName":"Scene1"} {}`,
		`{"action":"visibility","sceneName":"Scene0","sceneItemId":0}`, `{"action":"visibility","sceneName":"Scene0","enabled":false}`,
		`{"action":"visibility","sceneName":"Scene0","sceneItemId":-1,"enabled":false}`, `{"action":"visibility","sceneName":"Scene0","sceneItemId":0,"enabled":null}`,
		`{"action":"mute","inputName":"Camera"}`, `{"action":"mute","inputName":"Camera","muted":"false"}`,
		`{"action":"scene","sceneName":"Scene1","muted":false}`, `{"action":"scene","sceneName":"` + strings.Repeat("x", 4100) + `"}`,
	} {
		if rr := deckHTTP(h, "POST", body, "http://controller.test"); rr.Code != 400 {
			t.Errorf("invalid body status %d: %.120s", rr.Code, body)
		}
	}
	for _, origin := range []string{"", "null", "http://evil.test", "https://controller.test", "http://controller.test/path", "http://user@controller.test", "http://controller.test?x=1", "http://controller.test#x"} {
		if rr := deckHTTP(h, "POST", `{"action":"scene","sceneName":"Scene1"}`, origin); rr.Code != 403 {
			t.Errorf("origin %q: %d", origin, rr.Code)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) != 0 {
		t.Fatalf("invalid request reached OBS: %v", s.requests)
	}
}
func TestDeckNoOBS(t *testing.T) {
	h, err := NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	rr := deckHTTP(h, "GET", "", "")
	if rr.Code != 200 || !strings.Contains(rr.Body.String(), `"connected":false`) || !strings.Contains(rr.Body.String(), `"items":[]`) {
		t.Fatalf("offline GET %d %s", rr.Code, rr.Body)
	}
	if rr := deckHTTP(h, "POST", `{"action":"scene","sceneName":"Scene0"}`, "http://controller.test"); rr.Code != 503 {
		t.Fatalf("offline POST %d", rr.Code)
	}
}
func TestDeckBounds(t *testing.T) {
	s, _, h := deckFixture(t)
	s.scenes = 30
	s.items = 8
	rr := deckHTTP(h, "GET", "", "")
	if rr.Code != 200 {
		t.Fatal(rr.Code)
	}
	var got struct {
		Scenes []string `json:"scenes"`
		Items  []any    `json:"items"`
	}
	if json.Unmarshal(rr.Body.Bytes(), &got) != nil {
		t.Fatal(rr.Body.String())
	}
	if len(got.Scenes) != 24 || len(got.Items) != 128 {
		t.Fatalf("unbounded %d/%d", len(got.Scenes), len(got.Items))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, k := range s.requests {
		if k == "GetSceneItemList" {
			n++
		}
		if k == "GetGroupSceneItemList" {
			t.Fatal("nested traversal")
		}
	}
	if n != 16 {
		t.Fatalf("item queries %d", n)
	}
}
func TestDeckUnsupportedProbes(t *testing.T) {
	for _, kind := range []string{"GetInputMute", "GetMediaInputStatus"} {
		t.Run(kind, func(t *testing.T) {
			s, c, h := deckFixture(t)
			s.reject = kind
			rr := deckHTTP(h, "GET", "", "")
			if rr.Code != 200 {
				t.Fatalf("unsupported GET %d", rr.Code)
			}
			var got struct {
				Inputs []map[string]any `json:"inputs"`
			}
			json.Unmarshal(rr.Body.Bytes(), &got)
			field, body := "muted", `{"action":"mute","inputName":"Camera","muted":true}`
			if kind == "GetMediaInputStatus" {
				field = "mediaState"
				body = `{"action":"media-play","inputName":"Camera"}`
			}
			if _, ok := got.Inputs[0][field]; ok {
				t.Fatalf("unsupported field present %s", rr.Body)
			}
			rr = deckHTTP(h, "POST", body, "http://controller.test")
			if rr.Code != 409 {
				t.Fatalf("unsupported POST %d", rr.Code)
			}
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.conn == nil {
				t.Fatal("unsupported disconnected")
			}
		})
	}
}
func TestDeckFreshTargets(t *testing.T) {
	for _, body := range []string{`{"action":"scene","sceneName":"Missing"}`, `{"action":"visibility","sceneName":"Scene0","sceneItemId":99,"enabled":true,"targetId":"0000000000000000000000000000000000000000000000000000000000000000"}`, `{"action":"mute","inputName":"Missing","muted":true}`, `{"action":"media-play","inputName":"Missing"}`} {
		s, _, h := deckFixture(t)
		rr := deckHTTP(h, "POST", body, "http://controller.test")
		if rr.Code != 409 {
			t.Fatalf("target status %d", rr.Code)
		}
		s.mu.Lock()
		if s.mutated {
			t.Fatal("invalid target mutated")
		}
		s.mu.Unlock()
	}
	s, _, h := deckFixture(t)
	if deckHTTP(h, "GET", "", "").Code != 200 {
		t.Fatal("GET")
	}
	s.mu.Lock()
	s.scenes = 1
	s.mu.Unlock()
	if rr := deckHTTP(h, "POST", `{"action":"scene","sceneName":"Scene3"}`, "http://controller.test"); rr.Code != 409 {
		t.Fatalf("stale target %d", rr.Code)
	}
	// Existing first-three-scene API must not acquire Deck scene semantics.
	req := httptest.NewRequest("POST", "http://controller.test/api/scene", strings.NewReader(`{"sceneName":"Scene3"}`))
	req.Header.Set("Origin", "http://controller.test")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != 400 {
		t.Fatalf("legacy scene restriction %d", rr.Code)
	}
}
func TestDeckFailuresNeverReplay(t *testing.T) {
	for _, mode := range []string{"reject", "ignore", "disconnect-command", "disconnect-readback", "final-mismatch"} {
		t.Run(mode, func(t *testing.T) {
			s, c, h := deckFixture(t)
			want := 502
			switch mode {
			case "reject":
				s.reject = "SetCurrentProgramScene"
				want = 409
			case "ignore":
				s.ignore = true
			case "disconnect-command":
				s.disconnect = "SetCurrentProgramScene"
			case "disconnect-readback":
				s.disconnectAfter = true
			case "final-mismatch":
				s.revertFinal = true
			}
			rr := deckHTTP(h, "POST", `{"action":"scene","sceneName":"Scene3"}`, "http://controller.test")
			if rr.Code != want {
				t.Fatalf("failure = %d %s", rr.Code, rr.Body)
			}
			if strings.Contains(rr.Body.String(), "private-server-detail") || strings.Contains(rr.Body.String(), `"connected":true`) {
				t.Fatal("unsafe failure body")
			}
			c.mu.Lock()
			connected := c.conn != nil
			c.mu.Unlock()
			if strings.HasPrefix(mode, "disconnect") && connected {
				t.Fatal("transport failure kept connection")
			}
			s.mu.Lock()
			n := 0
			for _, k := range s.requests {
				if k == "SetCurrentProgramScene" {
					n++
				}
			}
			s.mu.Unlock()
			if n != 1 {
				t.Fatalf("mutation count %d", n)
			}
		})
	}
}
func TestDeckDiscoveryDisconnect(t *testing.T) {
	for _, kind := range []string{"GetSceneList", "GetSceneItemList", "GetInputList", "GetInputMute", "GetMediaInputStatus"} {
		t.Run(kind, func(t *testing.T) {
			s, c, h := deckFixture(t)
			s.disconnect = kind
			rr := deckHTTP(h, "GET", "", "")
			if rr.Code != 503 || strings.Contains(rr.Body.String(), `"connected":true`) {
				t.Fatalf("transport GET %d %s", rr.Code, rr.Body)
			}
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.conn != nil {
				t.Fatal("transport failure kept socket")
			}
		})
	}
}

func TestDeckActions(t *testing.T) {
	for _, tc := range []struct{ body, command, want string }{
		{`{"action":"scene","sceneName":"Scene3"}`, "SetCurrentProgramScene", "Scene3"},
		{`{"action":"visibility","sceneName":"Scene2","sceneItemId":0,"enabled":true}`, "SetSceneItemEnabled", "true"},
		{`{"action":"mute","inputName":"Camera","muted":true}`, "SetInputMute", "true"},
		{`{"action":"media-play","inputName":"Camera"}`, "TriggerMediaInputAction", "OBS_MEDIA_STATE_PLAYING"},
		{`{"action":"media-pause","inputName":"Camera"}`, "TriggerMediaInputAction", "OBS_MEDIA_STATE_PAUSED"},
		{`{"action":"media-restart","inputName":"Camera"}`, "TriggerMediaInputAction", "OBS_MEDIA_STATE_PLAYING"},
		{`{"action":"media-stop","inputName":"Camera"}`, "TriggerMediaInputAction", "OBS_MEDIA_STATE_STOPPED"},
	} {
		t.Run(tc.body, func(t *testing.T) {
			s, _, h := deckFixture(t)
			s.delayed = 2
			body := tc.body
			if tc.command == "SetSceneItemEnabled" {
				var inventory DeckState
				rr := deckHTTP(h, "GET", "", "")
				if err := json.Unmarshal(rr.Body.Bytes(), &inventory); err != nil {
					t.Fatal(err)
				}
				body = strings.TrimSuffix(body, "}") + fmt.Sprintf(`,"targetId":%q}`, inventory.Items[2].TargetID)
			}
			rr := deckHTTP(h, "POST", body, "http://controller.test")
			if rr.Code != 200 {
				t.Fatalf("POST = %d %s", rr.Code, rr.Body)
			}
			if !strings.Contains(rr.Body.String(), tc.want) {
				t.Fatalf("fresh state %s", rr.Body)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			n := 0
			for i, k := range s.requests {
				if k == tc.command {
					n++
					a := s.arguments[i]
					if k == "SetSceneItemEnabled" && (a["sceneName"] != "Scene2" || a["sceneItemId"] != float64(0) || a["sceneItemEnabled"] != true) {
						t.Fatalf("wrong target %v", a)
					}
				}
			}
			if n != 1 {
				t.Fatalf("mutation count = %d", n)
			}
		})
	}
}
