package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Explicit opt-in only: mutates the user's demo OBS, restores program scene,
// and removes the temporary scene/input even when an assertion fails.
func TestDeckLiveAcceptance(t *testing.T) {
	if os.Getenv("OBS_DECK_LIVE") != "1" {
		t.Skip("explicit demo OBS authorization required")
	}
	cfg, err := ReadConfiguration("../configuration.json")
	if err != nil {
		t.Fatal(err)
	}
	c := NewOBSClient(cfg)
	defer c.Close()
	request := func(k string, a any, v any) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.mu.Lock()
		defer c.mu.Unlock()
		if err := c.connect(ctx); err != nil {
			return err
		}
		return c.request(ctx, k, a, v)
	}
	var initial struct {
		Name string `json:"currentProgramSceneName"`
	}
	if err = request("GetCurrentProgramScene", map[string]any{}, &initial); err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("Deck acceptance %d", time.Now().UnixNano())
	input := name + " clip"
	sceneCreated, inputCreated := false, false
	defer func() {
		if err := request("SetCurrentProgramScene", map[string]any{"sceneName": initial.Name}, nil); err != nil {
			t.Errorf("restore scene: %v", err)
		}
		if inputCreated {
			if err := request("RemoveInput", map[string]any{"inputName": input}, nil); err != nil {
				t.Errorf("remove temporary input: %v", err)
			}
		}
		if sceneCreated {
			if err := request("RemoveScene", map[string]any{"sceneName": name}, nil); err != nil {
				t.Errorf("remove temporary scene: %v", err)
			}
		}
		clean := false
		for attempt := 0; attempt < 30; attempt++ {
			var scenes struct {
				Scenes []struct {
					Name string `json:"sceneName"`
				} `json:"scenes"`
			}
			var inputs struct {
				Inputs []struct {
					Name string `json:"inputName"`
				} `json:"inputs"`
			}
			if err := request("GetSceneList", map[string]any{}, &scenes); err != nil {
				t.Error(err)
				break
			}
			if err := request("GetInputList", map[string]any{}, &inputs); err != nil {
				t.Error(err)
				break
			}
			clean = true
			for _, scene := range scenes.Scenes {
				if scene.Name == name {
					clean = false
				}
			}
			for _, source := range inputs.Inputs {
				if source.Name == input {
					clean = false
				}
			}
			if clean {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if !clean {
			t.Error("temporary resources still exist after cleanup deadline")
		}
		var final struct {
			Name string `json:"currentProgramSceneName"`
		}
		if err := request("GetCurrentProgramScene", map[string]any{}, &final); err != nil || final.Name != initial.Name {
			t.Errorf("program restore not verified: %v", err)
		}
		if clean && final.Name == initial.Name {
			t.Log("temporary resources absent; initial program scene restored")
		}
	}()
	if err = request("CreateScene", map[string]any{"sceneName": name}, nil); err != nil {
		t.Fatal(err)
	}
	sceneCreated = true
	media, err := filepath.Abs("../test_output/deck-acceptance.mp4")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(media); err != nil {
		t.Fatal(err)
	}
	var created struct {
		ID int `json:"sceneItemId"`
	}
	if err = request("CreateInput", map[string]any{"sceneName": name, "inputName": input, "inputKind": "ffmpeg_source", "inputSettings": map[string]any{"local_file": media, "is_local_file": true, "looping": true, "restart_on_activate": true}, "sceneItemEnabled": true}, &created); err != nil {
		t.Fatal(err)
	}
	inputCreated = true
	h, err := newHandler(c)
	if err != nil {
		t.Fatal(err)
	}
	run := func(a map[string]any) DeckState {
		t.Helper()
		b, _ := json.Marshal(a)
		rr := deckHTTP(h, "POST", string(b), "http://controller.test")
		if rr.Code != 200 {
			t.Fatalf("live %s HTTP %d", a["action"], rr.Code)
		}
		var state DeckState
		if err := json.Unmarshal(rr.Body.Bytes(), &state); err != nil {
			t.Fatal(err)
		}
		t.Logf("live %s verified by HTTP and OBS readback", a["action"])
		return state
	}
	state := run(map[string]any{"action": "scene", "sceneName": name})
	var target DeckItem
	for _, i := range state.Items {
		if i.SceneName == name && i.SceneItemID == created.ID {
			target = i
		}
	}
	if target.TargetID == "" {
		t.Fatal("temporary source not discovered")
	}
	for _, enabled := range []bool{false, true} {
		run(map[string]any{"action": "visibility", "sceneName": name, "sceneItemId": created.ID, "targetId": target.TargetID, "enabled": enabled})
	}
	for _, muted := range []bool{true, false} {
		run(map[string]any{"action": "mute", "inputName": input, "muted": muted})
	}
	run(map[string]any{"action": "media-restart", "inputName": input})
	time.Sleep(1200 * time.Millisecond)
	run(map[string]any{"action": "media-pause", "inputName": input})
	run(map[string]any{"action": "media-play", "inputName": input})
	time.Sleep(1500 * time.Millisecond)
	run(map[string]any{"action": "media-restart", "inputName": input})
	var status struct {
		Cursor float64 `json:"mediaCursor"`
	}
	if err := request("GetMediaInputStatus", map[string]any{"inputName": input}, &status); err != nil || status.Cursor > 1000 {
		t.Fatalf("restart cursor=%v error=%v", status.Cursor, err)
	}
	run(map[string]any{"action": "media-stop", "inputName": input})
}
