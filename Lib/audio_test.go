package lib

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

type audioFixtureState struct {
	sync.Mutex
	volume                    float64
	muted, ignore, disconnect bool
	mutations                 int
}

func audioFixture(t *testing.T) (*audioFixtureState, *OBSClient, http.Handler) {
	t.Helper()
	s := &audioFixtureState{volume: .5}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		if c.WriteJSON(map[string]any{"op": 0, "d": map[string]any{}}) != nil {
			return
		}
		var hello any
		if c.ReadJSON(&hello) != nil {
			return
		}
		if c.WriteJSON(map[string]any{"op": 2, "d": map[string]any{}}) != nil {
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
			s.Lock()
			if s.disconnect {
				s.Unlock()
				return
			}
			data := map[string]any{}
			ok := true
			switch req.D.Kind {
			case "GetInputList":
				data["inputs"] = []any{map[string]any{"inputName": "Mic"}, map[string]any{"inputName": "Video"}}
			case "GetInputVolume":
				if req.D.Args["inputName"] == "Video" {
					ok = false
				} else {
					data["inputVolumeMul"] = s.volume
				}
			case "GetInputMute":
				data["inputMuted"] = s.muted
			case "SetInputVolume":
				s.mutations++
				if !s.ignore {
					s.volume = req.D.Args["inputVolumeMul"].(float64)
				}
			case "SetInputMute":
				s.mutations++
				if !s.ignore {
					s.muted = req.D.Args["inputMuted"].(bool)
				}
			default:
				ok = false
			}
			s.Unlock()
			if c.WriteJSON(map[string]any{"op": 7, "d": map[string]any{"requestId": req.D.ID, "requestStatus": map[string]any{"result": ok}, "responseData": data}}) != nil {
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

func audioRequest(h http.Handler, method, path, body, origin, content string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://localhost"+path, strings.NewReader(body))
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	if content != "" {
		r.Header.Set("Content-Type", content)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestAudioDiscoveryAndMutation(t *testing.T) {
	s, client, h := audioFixture(t)
	w := audioRequest(h, "GET", "/api/audio", "", "", "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var data struct {
		Inputs []audioInput `json:"inputs"`
	}
	if json.Unmarshal(w.Body.Bytes(), &data) != nil || len(data.Inputs) != 1 || data.Inputs[0].Name != "Mic" {
		t.Fatal(w.Body.String())
	}
	if client.conn == nil {
		t.Fatal("non-audio rejection disconnected client")
	}
	for _, tc := range []struct {
		body string
		code int
	}{
		{`{"action":"volume","inputName":"Mic","volume":0.7}`, 200},
		{`{"action":"mute","inputName":"Mic","muted":true}`, 200},
		{`{"action":"mute","inputName":"Video","muted":true}`, 409},
		{`{"action":"mute","inputName":"missing","muted":true}`, 409},
	} {
		w = audioRequest(h, "POST", "/api/audio/action", tc.body, "http://localhost", "application/json")
		if w.Code != tc.code {
			t.Errorf("%s: %d", tc.body, w.Code)
		}
	}
	s.Lock()
	s.ignore = true
	s.Unlock()
	w = audioRequest(h, "POST", "/api/audio/action", `{"action":"volume","inputName":"Mic","volume":0.2}`, "http://localhost", "application/json")
	if w.Code != 502 {
		t.Fatal("readback mismatch", w.Code)
	}
	s.Lock()
	s.volume = .9
	s.muted = false
	s.Unlock()
	w = audioRequest(h, "GET", "/api/audio", "", "", "")
	if !bytes.Contains(w.Body.Bytes(), []byte(`"volume":0.9`)) || !bytes.Contains(w.Body.Bytes(), []byte(`"muted":false`)) {
		t.Fatal("external update", w.Body.String())
	}
	s.Lock()
	s.disconnect = true
	s.Unlock()
	w = audioRequest(h, "GET", "/api/audio", "", "", "")
	if w.Code != 503 || client.conn != nil {
		t.Fatal("transport failure", w.Code)
	}
}

func TestAudioValidation(t *testing.T) {
	s, _, h := audioFixture(t)
	for _, body := range []string{`{}`, `{"action":"toggle","inputName":"Mic"}`, `{"action":"volume","inputName":"Mic","volume":-0.1}`, `{"action":"volume","inputName":"Mic","volume":1.1}`, `{"action":"volume","inputName":"Mic","volume":1e999}`, `{"action":"volume","inputName":"Mic","volume":null}`, `{"action":"mute","inputName":"Mic","muted":false,"extra":1}`, `{"action":"mute","inputName":"Mic","muted":false} {}`, `{"action":"mute","inputName":"Mic","muted":false,"muted":true}`} {
		w := audioRequest(h, "POST", "/api/audio/action", body, "http://localhost", "application/json")
		if w.Code != 400 {
			t.Errorf("%s: %d", body, w.Code)
		}
	}
	body := `{"action":"mute","inputName":"Mic","muted":true}`
	if w := audioRequest(h, "POST", "/api/audio/action", body, "http://evil", "application/json"); w.Code != 403 {
		t.Fatal(w.Code)
	}
	if w := audioRequest(h, "POST", "/api/audio/action", body, "http://localhost", "text/plain"); w.Code != 415 {
		t.Fatal(w.Code)
	}
	s.Lock()
	defer s.Unlock()
	if s.mutations != 0 {
		t.Fatal("invalid request mutated OBS")
	}
}
