package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestDeckRejectsReusedTarget(t *testing.T) {
	for _, change := range []string{"collection", "source"} {
		t.Run(change, func(t *testing.T) {
			s, _, h := deckFixture(t)
			rr := deckHTTP(h, "GET", "", "")
			var state map[string]any
			if err := json.Unmarshal(rr.Body.Bytes(), &state); err != nil {
				t.Fatal(err)
			}
			target := state["items"].([]any)[0].(map[string]any)["targetId"]
			if target == nil {
				t.Fatal("discovery missing target identity")
			}
			s.mu.Lock()
			if change == "collection" {
				s.collection = "Other"
			} else {
				s.sourceName = "Replacement"
			}
			s.mu.Unlock()
			body := fmt.Sprintf(`{"action":"visibility","sceneName":"Scene0","sceneItemId":0,"enabled":true,"targetId":%q}`, target)
			if rr = deckHTTP(h, "POST", body, "http://controller.test"); rr.Code != http.StatusConflict {
				t.Fatalf("stale target = %d", rr.Code)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			if s.mutated {
				t.Fatal("stale binding mutated OBS")
			}
		})
	}
}

func TestDeckRestartRequiresCursorReset(t *testing.T) {
	for _, ignore := range []bool{true, false} {
		t.Run(fmt.Sprint(ignore), func(t *testing.T) {
			s, _, h := deckFixture(t)
			s.mu.Lock()
			s.media = "OBS_MEDIA_STATE_PLAYING"
			s.cursor = 1000
			s.ignore = ignore
			s.mu.Unlock()
			rr := deckHTTP(h, "POST", `{"action":"media-restart","inputName":"Camera"}`, "http://controller.test")
			want := http.StatusOK
			if ignore {
				want = http.StatusBadGateway
			}
			if rr.Code != want {
				t.Fatalf("restart ignored=%v: got %d want %d", ignore, rr.Code, want)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			count := 0
			for _, k := range s.requests {
				if k == "TriggerMediaInputAction" {
					count++
				}
			}
			if count != 1 {
				t.Fatalf("sent %d commands", count)
			}
		})
	}
}
