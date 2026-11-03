package lib

import (
	"context"
	"os"
	"testing"
	"time"
)

// This test changes real OBS state. Run only with explicit authorization and
// exclusive control of recording; OBS has no recording ownership identifier.
func TestRecordLiveAcceptance(t *testing.T) {
	if os.Getenv("OBS_RECORD_LIVE") != "1" {
		t.Skip("explicit real OBS authorization required")
	}
	cfg, err := ReadConfiguration("../configuration.json")
	if err != nil {
		t.Fatal(err)
	}
	c := NewOBSClient(cfg)
	defer c.Close()
	snapshot := func() OBSStatus {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.snapshotLocked(ctx, false, true)
	}
	action := func(name string) (OBSStatus, int) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return c.Broadcast(ctx, name)
	}
	read := func(label string) OBSStatus {
		s := snapshot()
		if !s.Connected || !s.RecordAvailable {
			t.Fatalf("%s: fresh recording state unavailable", label)
		}
		t.Logf("%s: active=%v paused=%v duration=%.0f bytes=%.0f", label, s.RecordActive, s.RecordPaused, s.RecordDuration, s.RecordBytes)
		return s
	}
	initial := read("initial")
	if initial.RecordActive {
		t.Skip("recording already active; no commands sent")
	}
	owned := false
	defer func() {
		if !owned {
			return
		}
		s := snapshot()
		if !s.Connected || !s.RecordAvailable {
			t.Error("cleanup: fresh state unavailable; manual inspection required")
			return
		}
		if !s.RecordActive {
			t.Log("cleanup: fresh state inactive; no stop sent")
			return
		}
		_, code := action("stop-record")
		s = snapshot()
		if code != 200 || !s.Connected || !s.RecordAvailable || s.RecordActive {
			t.Errorf("cleanup: stop unconfirmed (HTTP %d); manual inspection required", code)
			return
		}
		t.Log("cleanup: stop confirmed by fresh read")
	}()
	s, code := action("start-record")
	if code != 200 || !s.Connected || !s.RecordAvailable || !s.RecordActive {
		t.Fatalf("start unconfirmed (HTTP %d); ownership uncertain, no cleanup stop; inspect OBS manually", code)
	}
	owned = true
	time.Sleep(5 * time.Second)
	s = read("after start")
	if !s.RecordActive {
		owned = false
		t.Fatal("recording stopped unexpectedly; no further commands")
	}
	if s.RecordPauseSupported == nil || !*s.RecordPauseSupported {
		t.Log("pause/resume unsupported or unknown; not verified")
	} else {
		for _, step := range []struct {
			name   string
			paused bool
		}{{"pause-record", true}, {"resume-record", false}} {
			_, code = action(step.name)
			s = read(step.name)
			if !s.RecordActive {
				owned = false
				t.Fatalf("%s: recording no longer active", step.name)
			}
			if code != 200 || s.RecordPaused != step.paused {
				t.Fatalf("%s: unconfirmed (HTTP %d)", step.name, code)
			}
			time.Sleep(3 * time.Second)
		}
	}
	s = read("before stop")
	if !s.RecordActive {
		owned = false
		t.Fatal("recording stopped unexpectedly; no stop sent")
	}
	_, code = action("stop-record")
	s = read("after stop")
	if !s.RecordActive {
		owned = false
	}
	if code != 200 || s.RecordActive {
		t.Fatalf("nominal stop unconfirmed (HTTP %d)", code)
	}
}
