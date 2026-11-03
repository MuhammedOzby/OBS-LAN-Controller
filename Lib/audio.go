package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type audioInput struct {
	Name   string  `json:"inputName"`
	Volume float64 `json:"volume"`
	Muted  bool    `json:"muted"`
}

type audioAction struct {
	Action    string   `json:"action"`
	InputName string   `json:"inputName"`
	Volume    *float64 `json:"volume"`
	Muted     *bool    `json:"muted"`
}

func (a audioAction) valid() bool {
	if a.InputName == "" {
		return false
	}
	switch a.Action {
	case "volume":
		return a.Volume != nil && a.Muted == nil && !math.IsNaN(*a.Volume) && !math.IsInf(*a.Volume, 0) && *a.Volume >= 0 && *a.Volume <= 1
	case "mute":
		return a.Muted != nil && a.Volume == nil
	}
	return false
}

func decodeAudioAction(r io.Reader) (audioAction, error) {
	var a audioAction
	d := json.NewDecoder(r)
	t, err := d.Token()
	if err != nil || t != json.Delim('{') {
		return a, errors.New("invalid object")
	}
	seen := map[string]bool{}
	for d.More() {
		t, err = d.Token()
		if err != nil {
			return a, err
		}
		key, ok := t.(string)
		if !ok || seen[key] {
			return a, errors.New("invalid field")
		}
		seen[key] = true
		var v json.RawMessage
		if err = d.Decode(&v); err != nil || string(v) == "null" {
			return a, errors.New("invalid value")
		}
		switch key {
		case "action":
			err = json.Unmarshal(v, &a.Action)
		case "inputName":
			err = json.Unmarshal(v, &a.InputName)
		case "volume":
			err = json.Unmarshal(v, &a.Volume)
		case "muted":
			err = json.Unmarshal(v, &a.Muted)
		default:
			return a, errors.New("unknown field")
		}
		if err != nil {
			return a, err
		}
	}
	if _, err = d.Token(); err != nil {
		return a, err
	}
	var extra any
	if d.Decode(&extra) != io.EOF || !a.valid() {
		return a, errors.New("invalid action")
	}
	return a, nil
}

// audioReadLocked probes both audio properties; ordinary OBS rejections mean
// the source has no audio (or disappeared), not a broken connection.
func (o *OBSClient) audioReadLocked(ctx context.Context, name string) (audioInput, error) {
	v := audioInput{Name: name}
	args := map[string]any{"inputName": name}
	var volume struct {
		Value *float64 `json:"inputVolumeMul"`
	}
	if err := o.deckRequest(ctx, "GetInputVolume", args, &volume); err != nil {
		return v, err
	}
	var mute struct {
		Value *bool `json:"inputMuted"`
	}
	if err := o.deckRequest(ctx, "GetInputMute", args, &mute); err != nil {
		return v, err
	}
	if volume.Value == nil || mute.Value == nil || math.IsNaN(*volume.Value) || math.IsInf(*volume.Value, 0) {
		o.disconnect()
		return v, errors.New("invalid audio response")
	}
	v.Volume, v.Muted = *volume.Value, *mute.Value
	return v, nil
}

func (o *OBSClient) audioListLocked(ctx context.Context) ([]audioInput, error) {
	var list struct {
		Inputs []struct {
			Name string `json:"inputName"`
		} `json:"inputs"`
	}
	if err := o.deckRequest(ctx, "GetInputList", map[string]any{}, &list); err != nil {
		return nil, err
	}
	result := make([]audioInput, 0)
	for i, input := range list.Inputs {
		if i >= deckMaxInputs {
			break
		}
		if input.Name == "" {
			continue
		}
		v, err := o.audioReadLocked(ctx, input.Name)
		if errors.Is(err, errOBSRejected) {
			continue
		}
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (o *OBSClient) audioChange(ctx context.Context, a audioAction) (audioInput, int) {
	var zero audioInput
	if !a.valid() {
		return zero, http.StatusBadRequest
	}
	if !o.lockDeck(ctx) {
		return zero, http.StatusServiceUnavailable
	}
	defer o.mu.Unlock()
	if o.connect(ctx) != nil {
		return zero, http.StatusServiceUnavailable
	}
	list, err := o.audioListLocked(ctx)
	if err != nil {
		return zero, http.StatusServiceUnavailable
	}
	found := false
	for _, v := range list {
		if v.Name == a.InputName {
			found = true
			break
		}
	}
	if !found {
		return zero, http.StatusConflict
	}
	args := map[string]any{"inputName": a.InputName}
	kind := "SetInputMute"
	if a.Action == "volume" {
		kind = "SetInputVolume"
		args["inputVolumeMul"] = *a.Volume
	} else {
		args["inputMuted"] = *a.Muted
	}
	if err = o.deckRequest(ctx, kind, args, nil); err != nil {
		if errors.Is(err, errOBSRejected) {
			return zero, http.StatusConflict
		}
		return zero, http.StatusBadGateway
	}
	v, err := o.audioReadLocked(ctx, a.InputName)
	if err != nil || (a.Action == "volume" && math.Abs(v.Volume-*a.Volume) > 0.0001) || (a.Action == "mute" && v.Muted != *a.Muted) {
		return zero, http.StatusBadGateway
	}
	return v, http.StatusOK
}

func audioRoutes(router *gin.Engine, obs *OBSClient) {
	router.GET("/api/audio", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if obs == nil {
			c.Status(503)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()
		if !obs.lockDeck(ctx) {
			c.Status(503)
			return
		}
		defer obs.mu.Unlock()
		if obs.connect(ctx) != nil {
			c.Status(503)
			return
		}
		list, err := obs.audioListLocked(ctx)
		if err != nil {
			c.Status(503)
			return
		}
		c.JSON(200, gin.H{"connected": true, "inputs": list})
	})
	router.POST("/api/audio/action", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		origin, err := url.Parse(c.GetHeader("Origin"))
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		if err != nil || origin.Scheme != scheme || origin.Host != c.Request.Host || origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
			c.Status(403)
			return
		}
		media, _, err := mime.ParseMediaType(c.GetHeader("Content-Type"))
		if err != nil || media != "application/json" {
			c.Status(415)
			return
		}
		_ = http.NewResponseController(c.Writer).SetReadDeadline(time.Now().Add(5 * time.Second))
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
		a, err := decodeAudioAction(c.Request.Body)
		if err != nil {
			c.Status(400)
			return
		}
		if obs == nil {
			c.Status(503)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
		defer cancel()
		v, code := obs.audioChange(ctx, a)
		if code != 200 {
			c.Status(code)
			return
		}
		c.JSON(200, v)
	})
}
