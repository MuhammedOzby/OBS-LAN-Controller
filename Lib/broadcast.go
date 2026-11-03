package lib

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type broadcastOutput struct {
	Active       *bool   `json:"outputActive"`
	Paused       bool    `json:"outputPaused"`
	Reconnecting bool    `json:"outputReconnecting"`
	Duration     float64 `json:"outputDuration"`
	Bytes        float64 `json:"outputBytes"`
}

type broadcastAction struct {
	command, status string
	active, paused  bool
}

func lookupBroadcastAction(action string) (broadcastAction, bool) {
	switch action {
	case "start-stream":
		return broadcastAction{"StartStream", "GetStreamStatus", true, false}, true
	case "stop-stream":
		return broadcastAction{"StopStream", "GetStreamStatus", false, false}, true
	case "start-record":
		return broadcastAction{"StartRecord", "GetRecordStatus", true, false}, true
	case "stop-record":
		return broadcastAction{"StopRecord", "GetRecordStatus", false, false}, true
	case "pause-record":
		return broadcastAction{"PauseRecord", "GetRecordStatus", true, true}, true
	case "resume-record":
		return broadcastAction{"ResumeRecord", "GetRecordStatus", true, false}, true
	}
	return broadcastAction{}, false
}

func (a broadcastAction) allowed(s broadcastOutput) bool {
	if s.Active == nil {
		return false
	}
	switch a.command {
	case "StartStream":
		return !*s.Active && !s.Reconnecting
	case "StopStream":
		return *s.Active || s.Reconnecting
	case "StartRecord":
		return !*s.Active
	case "StopRecord":
		return *s.Active
	case "PauseRecord":
		return *s.Active && !s.Paused
	case "ResumeRecord":
		return *s.Active && s.Paused
	}
	return false
}
func (a broadcastAction) matches(s broadcastOutput) bool {
	return s.Active != nil && *s.Active == a.active && (a.status != "GetRecordStatus" || !a.active || s.Paused == a.paused) && (a.status != "GetStreamStatus" || !s.Reconnecting)
}

// Broadcast executes once under the same mutex used for polling and scene changes.
// Preflight and readback never use cached display state. No failure retries a mutation.
func (o *OBSClient) Broadcast(ctx context.Context, action string) (OBSStatus, int) {
	a, ok := lookupBroadcastAction(action)
	if !ok {
		return OBSStatus{}, http.StatusBadRequest
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if err := o.connect(ctx); err != nil {
		return OBSStatus{}, http.StatusServiceUnavailable
	}
	var before broadcastOutput
	if err := o.request(ctx, a.status, map[string]any{}, &before); err != nil {
		if errors.Is(err, errOBSRejected) {
			return OBSStatus{}, http.StatusConflict
		}
		o.disconnect()
		return OBSStatus{}, http.StatusServiceUnavailable
	}
	if !a.allowed(before) {
		return OBSStatus{}, http.StatusConflict
	}
	if a.command == "PauseRecord" || a.command == "ResumeRecord" {
		if supported := o.recordPauseSupported(ctx); supported != nil && !*supported {
			return OBSStatus{}, http.StatusConflict
		}
	}
	o.sampled = time.Time{}
	if err := o.request(ctx, a.command, map[string]any{}, nil); err != nil {
		if errors.Is(err, errOBSRejected) {
			return OBSStatus{}, http.StatusConflict
		}
		o.disconnect()
		return OBSStatus{}, http.StatusBadGateway
	}
	for attempt := 0; attempt < 16; attempt++ {
		var after broadcastOutput
		if err := o.request(ctx, a.status, map[string]any{}, &after); err != nil {
			if !errors.Is(err, errOBSRejected) {
				o.disconnect()
			}
			return OBSStatus{}, http.StatusBadGateway
		}
		if a.matches(after) {
			status := o.snapshotLocked(ctx, false, true)
			// Verify the returned projection too: OBS may have changed externally.
			final := broadcastOutput{Active: &status.StreamActive, Reconnecting: status.Reconnecting}
			if a.status == "GetRecordStatus" {
				final = broadcastOutput{Active: &status.RecordActive, Paused: status.RecordPaused}
				if !status.RecordAvailable {
					return OBSStatus{}, http.StatusBadGateway
				}
			}
			if !status.Connected || !a.matches(final) {
				return OBSStatus{}, http.StatusBadGateway
			}
			return status, http.StatusOK
		}
		if attempt < 15 {
			select {
			case <-ctx.Done():
				return OBSStatus{}, http.StatusBadGateway
			case <-time.After(200 * time.Millisecond):
			}
		}
	}
	return OBSStatus{}, http.StatusBadGateway
}

func broadcastHandler(obs *OBSClient) gin.HandlerFunc {
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
		var request struct {
			Action    string `json:"action"`
			Confirmed bool   `json:"confirmed"`
		}
		decoder := json.NewDecoder(c.Request.Body)
		decoder.DisallowUnknownFields()
		if decoder.Decode(&request) != nil || !request.Confirmed {
			c.Status(http.StatusBadRequest)
			return
		}
		var extra any
		if decoder.Decode(&extra) != io.EOF {
			c.Status(http.StatusBadRequest)
			return
		}
		if _, ok := lookupBroadcastAction(request.Action); !ok {
			c.Status(http.StatusBadRequest)
			return
		}
		if obs == nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		status, code := obs.Broadcast(ctx, request.Action)
		if code != http.StatusOK {
			c.Status(code)
			return
		}
		c.JSON(http.StatusOK, status)
	}
}
