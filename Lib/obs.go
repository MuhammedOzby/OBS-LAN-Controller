package lib

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// OBSStatus is the browser-safe projection; configuration and credentials never leave the server.
type OBSStatus struct {
	Connected            bool     `json:"connected"`
	StreamActive         bool     `json:"streamActive"`
	StreamBytes          float64  `json:"streamBytes"`
	RecordAvailable      bool     `json:"recordAvailable"`
	RecordActive         bool     `json:"recordActive"`
	RecordPaused         bool     `json:"recordPaused"`
	RecordPauseSupported *bool    `json:"recordPauseSupported"`
	RecordDuration       float64  `json:"recordDuration"`
	RecordBytes          float64  `json:"recordBytes"`
	Reconnecting         bool     `json:"reconnecting"`
	ConnectionError      string   `json:"connectionError,omitempty"`
	Duration             float64  `json:"duration"`
	BitrateKbps          *float64 `json:"bitrateKbps"`
	CPU                  float64  `json:"cpu"`
	FPS                  float64  `json:"fps"`
	Memory               float64  `json:"memory"`
	RenderTime           float64  `json:"renderTime"`
	RenderSkipped        float64  `json:"renderSkipped"`
	RenderTotal          float64  `json:"renderTotal"`
	EncodeSkipped        float64  `json:"encodeSkipped"`
	EncodeTotal          float64  `json:"encodeTotal"`
	Dropped              float64  `json:"dropped"`
	Total                float64  `json:"total"`
	Congestion           float64  `json:"congestion"`
	Scenes               []string `json:"scenes"`
	CurrentScene         string   `json:"currentScene"`
	Preview              string   `json:"preview"`
	PreviewAvailable     bool     `json:"previewAvailable"`
}

type obsEnvelope struct {
	Op int             `json:"op"`
	D  json.RawMessage `json:"d"`
}

var errOBSRejected = errors.New("OBS isteği reddetti")
var errSceneNotAllowed = errors.New("geçersiz sahne")

// OBSClient serializes bounded requests and reconnects on the next poll, never replaying commands.
type OBSClient struct {
	mu             sync.Mutex
	config         Configuration
	conn           *websocket.Conn
	sequence       uint64
	cached         OBSStatus
	sampled        time.Time
	retryAfter     time.Time
	previousBytes  float64
	previousTime   time.Time
	previousActive bool
	favorites      *favoriteStore
}

func NewOBSClient(config Configuration) *OBSClient { return &OBSClient{config: config} }

func (o *OBSClient) Close() { o.mu.Lock(); defer o.mu.Unlock(); o.disconnect() }
func (o *OBSClient) disconnect() {
	if o.conn != nil {
		o.conn.Close()
		o.conn = nil
	}
	o.cached = OBSStatus{Scenes: []string{}}
	o.sampled = time.Time{}
	o.previousTime = time.Time{}
	o.previousActive = false
	o.retryAfter = time.Now().Add(2 * time.Second)
}

func (o *OBSClient) connect(ctx context.Context) error {
	if o.conn != nil {
		return nil
	}
	if time.Now().Before(o.retryAfter) {
		return errors.New("OBS yeniden bağlantısı bekleniyor")
	}
	dialer := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	conn, _, err := dialer.DialContext(ctx, "ws://"+net.JoinHostPort(o.config.OBS.IP, strconv.Itoa(o.config.OBS.Port)), nil)
	if err != nil {
		o.disconnect()
		return errors.New("OBS bağlantısı kurulamadı")
	}
	o.conn = conn
	conn.SetReadLimit(4 << 20)
	deadline := time.Now().Add(3 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	conn.SetReadDeadline(deadline)
	conn.SetWriteDeadline(deadline)
	var hello obsEnvelope
	if err = conn.ReadJSON(&hello); err != nil || hello.Op != 0 {
		o.disconnect()
		return errors.New("OBS karşılama hatası")
	}
	var data struct {
		Authentication *struct {
			Challenge string `json:"challenge"`
			Salt      string `json:"salt"`
		} `json:"authentication"`
	}
	if json.Unmarshal(hello.D, &data) != nil {
		o.disconnect()
		return errors.New("OBS protokol hatası")
	}
	identify := map[string]any{"rpcVersion": 1, "eventSubscriptions": 0}
	if data.Authentication != nil {
		secret := sha256.Sum256([]byte(o.config.OBS.Password + data.Authentication.Salt))
		auth := sha256.Sum256([]byte(base64.StdEncoding.EncodeToString(secret[:]) + data.Authentication.Challenge))
		identify["authentication"] = base64.StdEncoding.EncodeToString(auth[:])
	}
	if err = conn.WriteJSON(map[string]any{"op": 1, "d": identify}); err != nil {
		o.disconnect()
		return errors.New("OBS tanımlama hatası")
	}
	var identified obsEnvelope
	if err = conn.ReadJSON(&identified); err != nil || identified.Op != 2 {
		o.disconnect()
		return errors.New("OBS kimlik doğrulaması başarısız")
	}
	return nil
}

func (o *OBSClient) request(ctx context.Context, kind string, args any, target any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	deadline := time.Now().Add(3 * time.Second)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	o.conn.SetReadDeadline(deadline)
	o.conn.SetWriteDeadline(deadline)
	o.sequence++
	id := strconv.FormatUint(o.sequence, 10)
	if err := o.conn.WriteJSON(map[string]any{"op": 6, "d": map[string]any{"requestType": kind, "requestId": id, "requestData": args}}); err != nil {
		return errors.New("OBS isteği gönderilemedi")
	}
	for {
		var envelope obsEnvelope
		if err := o.conn.ReadJSON(&envelope); err != nil {
			return errors.New("OBS yanıtı alınamadı")
		}
		if envelope.Op != 7 {
			continue
		}
		var response struct {
			RequestID     string `json:"requestId"`
			RequestStatus struct {
				Result bool `json:"result"`
			} `json:"requestStatus"`
			ResponseData json.RawMessage `json:"responseData"`
		}
		if err := json.Unmarshal(envelope.D, &response); err != nil {
			return errors.New("OBS yanıtı geçersiz")
		}
		if response.RequestID != id {
			continue
		}
		if !response.RequestStatus.Result {
			return errOBSRejected
		}
		if target != nil {
			if err := json.Unmarshal(response.ResponseData, target); err != nil {
				return errors.New("OBS verisi geçersiz")
			}
		}
		return nil
	}
}

func bitrate(active bool, previous, current float64, before, now time.Time) *float64 {
	seconds := now.Sub(before).Seconds()
	if !active || before.IsZero() || seconds <= 0 || current < previous {
		return nil
	}
	result := (current - previous) * 8 / seconds / 1000
	return &result
}

func (o *OBSClient) Snapshot(ctx context.Context, preview bool) OBSStatus {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.snapshotLocked(ctx, preview, false)
}

// snapshotLocked requires mu; fresh bypasses the display-only one-second cache.
func (o *OBSClient) snapshotLocked(ctx context.Context, preview, fresh bool) OBSStatus {
	if ctx.Err() != nil {
		return OBSStatus{Scenes: []string{}}
	}
	if !fresh && time.Since(o.sampled) < time.Second {
		s := o.cached
		if !preview {
			s.Preview = ""
			s.PreviewAvailable = false
		}
		return s
	}
	if err := o.connect(ctx); err != nil {
		s := o.cached
		// Only fixed, locally defined messages are exposed. Never expose a dial,
		// WebSocket, or OBS response error which could contain configuration.
		switch err.Error() {
		case "OBS kimlik doğrulaması başarısız":
			s.ConnectionError = "OBS kimlik doğrulaması başarısız; bağlantı ayarlarını kontrol edin."
		case "OBS yeniden bağlantısı bekleniyor":
			s.ConnectionError = "OBS bağlantısı yeniden deneniyor."
		default:
			s.ConnectionError = "OBS bağlantısı kurulamadı; OBS WebSocket ayarlarını kontrol edin."
		}
		return s
	}
	var scenes struct {
		Current string `json:"currentProgramSceneName"`
		Scenes  []struct {
			Name string `json:"sceneName"`
		} `json:"scenes"`
	}
	var stats struct {
		CPU           float64 `json:"cpuUsage"`
		FPS           float64 `json:"activeFps"`
		Memory        float64 `json:"memoryUsage"`
		Render        float64 `json:"averageFrameRenderTime"`
		RenderSkipped float64 `json:"renderSkippedFrames"`
		RenderTotal   float64 `json:"renderTotalFrames"`
		EncodeSkipped float64 `json:"outputSkippedFrames"`
		EncodeTotal   float64 `json:"outputTotalFrames"`
	}
	var stream struct {
		Active       bool    `json:"outputActive"`
		Reconnecting bool    `json:"outputReconnecting"`
		Duration     float64 `json:"outputDuration"`
		Bytes        float64 `json:"outputBytes"`
		Skipped      float64 `json:"outputSkippedFrames"`
		Total        float64 `json:"outputTotalFrames"`
		Congestion   float64 `json:"outputCongestion"`
	}
	for _, r := range []struct {
		kind   string
		target any
	}{{"GetSceneList", &scenes}, {"GetStats", &stats}, {"GetStreamStatus", &stream}} {
		if err := o.request(ctx, r.kind, map[string]any{}, r.target); err != nil {
			o.disconnect()
			return o.cached
		}
	}
	var record broadcastOutput
	recordErr := o.request(ctx, "GetRecordStatus", map[string]any{}, &record)
	if recordErr != nil && !errors.Is(recordErr, errOBSRejected) {
		o.disconnect()
		return o.cached
	}
	now := time.Now()
	s := OBSStatus{Connected: true, StreamActive: stream.Active, StreamBytes: stream.Bytes, Reconnecting: stream.Reconnecting, Duration: stream.Duration, CPU: stats.CPU, FPS: stats.FPS, Memory: stats.Memory, RenderTime: stats.Render, RenderSkipped: stats.RenderSkipped, RenderTotal: stats.RenderTotal, EncodeSkipped: stats.EncodeSkipped, EncodeTotal: stats.EncodeTotal, Dropped: stream.Skipped, Total: stream.Total, Congestion: stream.Congestion, CurrentScene: scenes.Current, Scenes: []string{}}
	if recordErr == nil && record.Active != nil {
		s.RecordAvailable = true
		s.RecordActive = *record.Active
		s.RecordPaused = record.Paused
		s.RecordDuration = record.Duration
		s.RecordBytes = record.Bytes
		s.RecordPauseSupported = o.recordPauseSupported(ctx)
	}
	for _, scene := range scenes.Scenes {
		s.Scenes = append(s.Scenes, scene.Name)
	}
	if o.previousActive {
		s.BitrateKbps = bitrate(stream.Active && !stream.Reconnecting, o.previousBytes, stream.Bytes, o.previousTime, now)
	}
	o.previousBytes = stream.Bytes
	o.previousTime = now
	o.previousActive = stream.Active
	if preview && scenes.Current != "" {
		var image struct {
			Data string `json:"imageData"`
		}
		err := o.request(ctx, "GetSourceScreenshot", map[string]any{"sourceName": scenes.Current, "imageFormat": "jpeg", "imageWidth": 960, "imageCompressionQuality": 65}, &image)
		if err == nil {
			s.Preview = image.Data
			s.PreviewAvailable = image.Data != ""
		} else {
			if !errors.Is(err, errOBSRejected) {
				o.disconnect()
				return o.cached
			}
			s.PreviewAvailable = false
		}
	}
	o.cached = s
	o.sampled = now
	return s
}

func (o *OBSClient) SetScene(ctx context.Context, name string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	// Scene authorization is based on a fresh OBS read under the same lock as
	// the command, never on the dashboard's one-second display cache.
	status := o.snapshotLocked(ctx, false, true)
	if !status.Connected || o.conn == nil {
		return errors.New("OBS bağlı değil")
	}
	var config *favoriteConfig
	if o.favorites != nil {
		config = o.favorites.get()
	}
	collection := ""
	if config != nil {
		var err error
		collection, err = o.collectionLocked(ctx)
		if err != nil {
			return err
		}
	}
	allowed := false
	for _, scene := range favoriteSlots(config, collection, status.Scenes) {
		if scene == name && containsScene(status.Scenes, name) {
			allowed = true
		}
	}
	if !allowed {
		return errSceneNotAllowed
	}
	if err := o.request(ctx, "SetCurrentProgramScene", map[string]any{"sceneName": name}, nil); err != nil {
		o.disconnect()
		return errors.New("sahne değiştirilemedi")
	}
	for attempt := 0; attempt < 5; attempt++ {
		var current struct {
			Name string `json:"currentProgramSceneName"`
		}
		if err := o.request(ctx, "GetCurrentProgramScene", map[string]any{}, &current); err != nil {
			o.disconnect()
			return errors.New("sahne doğrulanamadı")
		}
		if current.Name == name {
			o.cached.CurrentScene = name
			o.sampled = time.Time{}
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	o.sampled = time.Time{}
	return errors.New("sahne değişimi henüz doğrulanamadı")
}
