package lib

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	deckMaxScenes = 24
	deckMaxItems  = 128
	deckMaxInputs = 128
)

type DeckState struct {
	Connected    bool        `json:"connected"`
	CurrentScene string      `json:"currentScene"`
	Scenes       []string    `json:"scenes"`
	Items        []DeckItem  `json:"items"`
	Inputs       []DeckInput `json:"inputs"`
}
type DeckItem struct {
	TargetID    string `json:"targetId"`
	SceneName   string `json:"sceneName"`
	SceneItemID int    `json:"sceneItemId"`
	SourceName  string `json:"sourceName"`
	Enabled     bool   `json:"enabled"`
}
type DeckInput struct {
	InputName  string `json:"inputName"`
	InputKind  string `json:"inputKind"`
	Muted      *bool  `json:"muted,omitempty"`
	MediaState string `json:"mediaState,omitempty"`
}

func emptyDeck() DeckState {
	return DeckState{Scenes: []string{}, Items: []DeckItem{}, Inputs: []DeckInput{}}
}

// A context-aware mutex acquisition keeps the HTTP deadline bounded even when
// the ordinary dashboard poll owns the shared OBS connection.
func (o *OBSClient) lockDeck(ctx context.Context) bool {
	for ctx.Err() == nil {
		if o.mu.TryLock() {
			return true
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(10 * time.Millisecond):
		}
	}
	return false
}

// deckRequest must run under mu. Unsupported requests are ordinary OBS rejections;
// transport/protocol failures invalidate the connection, never yielding stale state.
func (o *OBSClient) deckRequest(ctx context.Context, kind string, args, target any) error {
	err := o.request(ctx, kind, args, target)
	if err != nil && !errors.Is(err, errOBSRejected) {
		o.disconnect()
	}
	return err
}

// deckLocked discovers only direct scene items. Groups are exposed as direct
// items but never traversed, so nested sources are intentionally not assignable.
// A global item/input cap also bounds probe requests and response sizes.
func (o *OBSClient) deckLocked(ctx context.Context) (DeckState, error) {
	s := emptyDeck()
	var collection struct {
		Name string `json:"currentSceneCollectionName"`
	}
	if err := o.deckRequest(ctx, "GetSceneCollectionList", map[string]any{}, &collection); err != nil || collection.Name == "" {
		return emptyDeck(), errors.New("OBS sahne koleksiyonu okunamadı")
	}
	var scenes struct {
		Current string `json:"currentProgramSceneName"`
		Scenes  []struct {
			Name string `json:"sceneName"`
			UUID string `json:"sceneUuid"`
		} `json:"scenes"`
	}
	if err := o.deckRequest(ctx, "GetSceneList", map[string]any{}, &scenes); err != nil {
		return emptyDeck(), err
	}
	s.CurrentScene = scenes.Current
	for i, scene := range scenes.Scenes {
		if i >= deckMaxScenes {
			break
		}
		s.Scenes = append(s.Scenes, scene.Name)
		if len(s.Items) >= deckMaxItems {
			continue
		}
		var items struct {
			Items []struct {
				ID      int    `json:"sceneItemId"`
				Name    string `json:"sourceName"`
				UUID    string `json:"sourceUuid"`
				Enabled bool   `json:"sceneItemEnabled"`
			} `json:"sceneItems"`
		}
		if err := o.deckRequest(ctx, "GetSceneItemList", map[string]any{"sceneName": scene.Name}, &items); err != nil {
			return emptyDeck(), err
		}
		for _, item := range items.Items {
			if len(s.Items) >= deckMaxItems {
				break
			}
			identity, _ := json.Marshal([]any{collection.Name, scene.Name, scene.UUID, item.ID, item.Name, item.UUID})
			s.Items = append(s.Items, DeckItem{TargetID: fmt.Sprintf("%x", sha256.Sum256(identity)), SceneName: scene.Name, SceneItemID: item.ID, SourceName: item.Name, Enabled: item.Enabled})
		}
	}
	var inputs struct {
		Inputs []DeckInput `json:"inputs"`
	}
	if err := o.deckRequest(ctx, "GetInputList", map[string]any{}, &inputs); err != nil {
		return emptyDeck(), err
	}
	for i, input := range inputs.Inputs {
		if i >= deckMaxInputs {
			break
		}
		// Capabilities must come from their dedicated OBS probes, not the list.
		input.Muted = nil
		input.MediaState = ""
		var mute struct {
			Muted *bool `json:"inputMuted"`
		}
		err := o.deckRequest(ctx, "GetInputMute", map[string]any{"inputName": input.InputName}, &mute)
		if err == nil {
			input.Muted = mute.Muted
		} else if !errors.Is(err, errOBSRejected) {
			return emptyDeck(), err
		}
		var media struct {
			State string `json:"mediaState"`
		}
		err = o.deckRequest(ctx, "GetMediaInputStatus", map[string]any{"inputName": input.InputName}, &media)
		if err == nil {
			input.MediaState = media.State
		} else if !errors.Is(err, errOBSRejected) {
			return emptyDeck(), err
		}
		s.Inputs = append(s.Inputs, input)
	}
	s.Connected = true
	return s, nil
}

// DeckAction deliberately has no arbitrary OBS request, macro or output action.
// Pointer fields distinguish explicit false/zero from a missing target/value.
type DeckAction struct {
	TargetID    string `json:"targetId,omitempty"`
	Action      string `json:"action"`
	SceneName   string `json:"sceneName,omitempty"`
	SceneItemID *int   `json:"sceneItemId,omitempty"`
	InputName   string `json:"inputName,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Muted       *bool  `json:"muted,omitempty"`
}

func (a DeckAction) valid() bool {
	if a.Action != "visibility" && a.TargetID != "" {
		return false
	}
	switch a.Action {
	case "scene":
		return a.SceneName != "" && a.SceneItemID == nil && a.InputName == "" && a.Enabled == nil && a.Muted == nil
	case "visibility":
		return a.TargetID != "" && len(a.TargetID) == 64 && a.SceneName != "" && a.SceneItemID != nil && *a.SceneItemID >= 0 && a.Enabled != nil && a.InputName == "" && a.Muted == nil
	case "mute":
		return a.InputName != "" && a.Muted != nil && a.SceneName == "" && a.SceneItemID == nil && a.Enabled == nil
	case "media-play", "media-pause", "media-restart", "media-stop":
		return a.InputName != "" && a.SceneName == "" && a.SceneItemID == nil && a.Enabled == nil && a.Muted == nil
	}
	return false
}
func (a DeckAction) mediaTarget() (string, string) {
	switch a.Action {
	case "media-play":
		return "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PLAY", "OBS_MEDIA_STATE_PLAYING"
	case "media-pause":
		return "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_PAUSE", "OBS_MEDIA_STATE_PAUSED"
	case "media-restart":
		return "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_RESTART", "OBS_MEDIA_STATE_PLAYING"
	case "media-stop":
		return "OBS_WEBSOCKET_MEDIA_INPUT_ACTION_STOP", "OBS_MEDIA_STATE_STOPPED"
	}
	return "", ""
}
func (a DeckAction) targetExists(s DeckState) bool {
	switch a.Action {
	case "scene":
		for _, name := range s.Scenes {
			if name == a.SceneName {
				return true
			}
		}
	case "visibility":
		for _, item := range s.Items {
			if item.SceneName == a.SceneName && item.SceneItemID == *a.SceneItemID && item.TargetID == a.TargetID {
				return true
			}
		}
	default:
		for _, input := range s.Inputs {
			if input.InputName == a.InputName {
				if a.Action == "mute" {
					return input.Muted != nil
				}
				return input.MediaState != "" && input.MediaState != "OBS_MEDIA_STATE_NONE" && input.MediaState != "OBS_MEDIA_STATE_ERROR"
			}
		}
	}
	return false
}
func (a DeckAction) matches(s DeckState) bool {
	switch a.Action {
	case "scene":
		return s.CurrentScene == a.SceneName
	case "visibility":
		for _, item := range s.Items {
			if item.SceneName == a.SceneName && item.SceneItemID == *a.SceneItemID && item.TargetID == a.TargetID {
				return item.Enabled == *a.Enabled
			}
		}
	case "mute":
		for _, input := range s.Inputs {
			if input.InputName == a.InputName {
				return input.Muted != nil && *input.Muted == *a.Muted
			}
		}
	default:
		_, state := a.mediaTarget()
		for _, input := range s.Inputs {
			if input.InputName == a.InputName {
				return input.MediaState == state
			}
		}
	}
	return false
}
func (a DeckAction) command() (string, map[string]any) {
	switch a.Action {
	case "scene":
		return "SetCurrentProgramScene", map[string]any{"sceneName": a.SceneName}
	case "visibility":
		return "SetSceneItemEnabled", map[string]any{"sceneName": a.SceneName, "sceneItemId": *a.SceneItemID, "sceneItemEnabled": *a.Enabled}
	case "mute":
		return "SetInputMute", map[string]any{"inputName": a.InputName, "inputMuted": *a.Muted}
	default:
		action, _ := a.mediaTarget()
		return "TriggerMediaInputAction", map[string]any{"inputName": a.InputName, "mediaAction": action}
	}
}

// readDeckTarget asks OBS about the exact requested target; no display cache or
// optimistic mutation result is ever used as verification.
func (o *OBSClient) readDeckTarget(ctx context.Context, a DeckAction) (bool, error) {
	switch a.Action {
	case "scene":
		var v struct {
			Name *string `json:"currentProgramSceneName"`
		}
		err := o.deckRequest(ctx, "GetCurrentProgramScene", map[string]any{}, &v)
		return v.Name != nil && *v.Name == a.SceneName, err
	case "visibility":
		var v struct {
			Enabled *bool `json:"sceneItemEnabled"`
		}
		err := o.deckRequest(ctx, "GetSceneItemEnabled", map[string]any{"sceneName": a.SceneName, "sceneItemId": *a.SceneItemID}, &v)
		return v.Enabled != nil && *v.Enabled == *a.Enabled, err
	case "mute":
		var v struct {
			Muted *bool `json:"inputMuted"`
		}
		err := o.deckRequest(ctx, "GetInputMute", map[string]any{"inputName": a.InputName}, &v)
		return v.Muted != nil && *v.Muted == *a.Muted, err
	default:
		var v struct {
			State string `json:"mediaState"`
		}
		err := o.deckRequest(ctx, "GetMediaInputStatus", map[string]any{"inputName": a.InputName}, &v)
		_, state := a.mediaTarget()
		return v.State == state, err
	}
}
func (o *OBSClient) deckAction(ctx context.Context, a DeckAction) (DeckState, int) {
	if !a.valid() {
		return emptyDeck(), http.StatusBadRequest
	}
	if !o.lockDeck(ctx) {
		return emptyDeck(), http.StatusServiceUnavailable
	}
	defer o.mu.Unlock()
	if o.connect(ctx) != nil {
		return emptyDeck(), http.StatusServiceUnavailable
	}
	before, err := o.deckLocked(ctx)
	if err != nil {
		return emptyDeck(), http.StatusServiceUnavailable
	}
	if !a.targetExists(before) {
		return emptyDeck(), http.StatusConflict
	}
	kind, args := a.command()
	var cursorBefore *float64
	if a.Action == "media-restart" {
		var media struct {
			Cursor *float64 `json:"mediaCursor"`
		}
		if err := o.deckRequest(ctx, "GetMediaInputStatus", map[string]any{"inputName": a.InputName}, &media); err != nil || media.Cursor == nil {
			return emptyDeck(), http.StatusBadGateway
		}
		cursorBefore = media.Cursor
	}
	commandAt := time.Now()
	o.sampled = time.Time{} // Invalidate the unrelated dashboard cache, not its scene policy.
	if err = o.deckRequest(ctx, kind, args, nil); err != nil {
		if errors.Is(err, errOBSRejected) {
			return emptyDeck(), http.StatusConflict
		}
		return emptyDeck(), http.StatusBadGateway
	}
	// This subdeadline includes the actual read requests, not just the sleeps.
	// It leaves room in the eight-second HTTP budget for fresh final discovery.
	verify, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	for attempt := 0; attempt < 16; attempt++ {
		matched, err := o.readDeckTarget(verify, a)
		// PLAYING alone does not prove a restart. Verify the cursor is near
		// the beginning, allowing only time elapsed since this command.
		if err == nil && matched && a.Action == "media-restart" {
			var media struct {
				Cursor *float64 `json:"mediaCursor"`
			}
			err = o.deckRequest(verify, "GetMediaInputStatus", map[string]any{"inputName": a.InputName}, &media)
			matched = media.Cursor != nil && *media.Cursor >= 0 && *media.Cursor <= float64(time.Since(commandAt).Milliseconds())+250 && (*cursorBefore <= 250 || *media.Cursor < *cursorBefore)
		}
		if err != nil {
			return emptyDeck(), http.StatusBadGateway
		}
		if matched {
			after, err := o.deckLocked(ctx)
			if err != nil || !a.targetExists(after) || !a.matches(after) {
				return emptyDeck(), http.StatusBadGateway
			}
			return after, http.StatusOK
		}
		if attempt < 15 {
			select {
			case <-verify.Done():
				return emptyDeck(), http.StatusBadGateway
			case <-time.After(150 * time.Millisecond):
			}
		}
	}
	return emptyDeck(), http.StatusBadGateway
}

// decodeDeckAction rejects duplicate, case-variant and explicit null fields as
// well as unknown fields. encoding/json alone accepts all three ambiguities.
func decodeDeckAction(r io.Reader) (DeckAction, error) {
	var a DeckAction
	d := json.NewDecoder(r)
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return a, errors.New("invalid action")
	}
	seen := map[string]bool{}
	for d.More() {
		token, err = d.Token()
		if err != nil {
			return a, err
		}
		key, ok := token.(string)
		if !ok || seen[key] {
			return a, errors.New("invalid field")
		}
		seen[key] = true
		switch key {
		case "action", "sceneName", "sceneItemId", "inputName", "enabled", "muted", "targetId":
		default:
			return a, errors.New("invalid field")
		}
		var value json.RawMessage
		if err = d.Decode(&value); err != nil {
			return a, err
		}
		if string(value) == "null" {
			return a, errors.New("invalid value")
		}
		switch key {
		case "targetId":
			err = json.Unmarshal(value, &a.TargetID)
		case "action":
			err = json.Unmarshal(value, &a.Action)
		case "sceneName":
			err = json.Unmarshal(value, &a.SceneName)
		case "sceneItemId":
			err = json.Unmarshal(value, &a.SceneItemID)
		case "inputName":
			err = json.Unmarshal(value, &a.InputName)
		case "enabled":
			err = json.Unmarshal(value, &a.Enabled)
		case "muted":
			err = json.Unmarshal(value, &a.Muted)
		}
		if err != nil {
			return a, err
		}
	}
	if _, err = d.Token(); err != nil {
		return a, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return a, errors.New("trailing data")
	}
	if !a.valid() {
		return a, errors.New("invalid action")
	}
	return a, nil
}

func deckActionHandler(obs *OBSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		origin, err := url.Parse(c.GetHeader("Origin"))
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		if err != nil || origin.Scheme != scheme || origin.Host != c.Request.Host || origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
			c.Status(http.StatusForbidden)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
		a, err := decodeDeckAction(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if obs == nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()
		s, code := obs.deckAction(ctx, a)
		if code != http.StatusOK {
			c.Status(code)
			return
		}
		c.JSON(http.StatusOK, s)
	}
}

func deckHandler(obs *OBSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if obs == nil {
			c.JSON(http.StatusOK, emptyDeck())
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()
		if !obs.lockDeck(ctx) {
			c.JSON(http.StatusServiceUnavailable, emptyDeck())
			return
		}
		defer obs.mu.Unlock()
		if obs.connect(ctx) != nil {
			c.JSON(http.StatusServiceUnavailable, emptyDeck())
			return
		}
		s, err := obs.deckLocked(ctx)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, emptyDeck())
			return
		}
		c.JSON(http.StatusOK, s)
	}
}
