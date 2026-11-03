package lib

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

const previewMaxViewers = 3

type PreviewStatus struct {
	Connected bool    `json:"connected"`
	Viewers   int     `json:"viewers"`
	Streaming bool    `json:"streaming"`
	Width     uint32  `json:"width"`
	Height    uint32  `json:"height"`
	FPS       float64 `json:"fps"`
	Error     string  `json:"error"`
}
type previewViewer struct {
	id             string
	pc             *webrtc.PeerConnection
	track          *webrtc.TrackLocalStaticSample
	frames         chan media.Sample
	done           chan struct{}
	ready, waiting bool // protected by service mutex
}
type previewService struct {
	mu          sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
	listener    *net.UnixListener
	cleanup     func()
	wg          sync.WaitGroup
	source      *net.UnixConn
	viewers     map[string]*previewViewer
	wake        chan struct{}
	config      previewConfig
	generation  uint32
	lastKey     time.Time
	lastKeySent time.Time
	lastFrame   time.Time
	lastPTS     int64
	havePTS     bool
	errorText   string
	authorize   func(*http.Request) bool // set before HTTP server starts
}

func newPreviewService() *previewService {
	ctx, cancel := context.WithCancel(context.Background())
	return &previewService{ctx: ctx, cancel: cancel, viewers: make(map[string]*previewViewer), wake: make(chan struct{}, 1)}
}
func (p *previewService) start() error {
	listener, cleanup, err := listenPreview()
	if err != nil {
		p.mu.Lock()
		p.errorText = err.Error()
		p.mu.Unlock()
		return err
	}
	p.listener, p.cleanup = listener, cleanup
	p.wg.Add(1)
	go func() { defer p.wg.Done(); p.accept() }()
	return nil
}
func (p *previewService) close() {
	p.cancel()
	if p.listener != nil {
		p.listener.Close()
	}
	p.mu.Lock()
	if p.source != nil {
		p.source.Close()
	}
	viewers := make([]*previewViewer, 0, len(p.viewers))
	for _, v := range p.viewers {
		viewers = append(viewers, v)
	}
	p.mu.Unlock()
	for _, v := range viewers {
		p.removeViewer(v.id)
	}
	p.wg.Wait()
	if p.cleanup != nil {
		p.cleanup()
	}
}
func (p *previewService) notify() {
	select {
	case p.wake <- struct{}{}:
	default:
	}
}
func (p *previewService) status() PreviewStatus {
	p.mu.Lock()
	defer p.mu.Unlock()
	status := PreviewStatus{Connected: p.source != nil, Viewers: p.activeLocked(), Error: p.errorText}
	status.Streaming = status.Connected && status.Viewers > 0 && !p.lastFrame.IsZero() && time.Since(p.lastFrame) < 3*time.Second
	if status.Streaming {
		status.Width = p.config.width
		status.Height = p.config.height
		status.FPS = float64(p.config.fpsNum) / float64(p.config.fpsDen)
	}
	return status
}
func (p *previewService) activeLocked() int {
	n := 0
	for _, v := range p.viewers {
		if v.ready {
			n++
		}
	}
	return n
}
func (p *previewService) accept() {
	for {
		conn, err := p.listener.AcceptUnix()
		if err != nil {
			return
		}
		if !previewPeerAllowed(conn) {
			conn.Close()
			continue
		}
		p.mu.Lock()
		if p.source != nil {
			p.mu.Unlock()
			conn.Close()
			continue
		}
		// Reserve one source before handshake; malformed peers cannot spawn workers.
		p.source = conn
		p.mu.Unlock()
		p.wg.Add(1)
		go func() { defer p.wg.Done(); p.serveSource(conn) }()
	}
}
func (p *previewService) serveSource(conn *net.UnixConn) {
	defer func() {
		conn.Close()
		p.mu.Lock()
		p.source = nil
		p.config = previewConfig{}
		p.generation = 0
		p.lastFrame = time.Time{}
		p.havePTS = false
		for _, v := range p.viewers {
			v.waiting = true
			drainSamples(v.frames)
		}
		p.mu.Unlock()
	}()
	conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	var hello [8]byte
	if _, err := io.ReadFull(conn, hello[:]); err != nil || string(hello[:]) != previewMagic {
		return
	}
	if _, err := io.WriteString(conn, previewMagic); err != nil {
		return
	}
	conn.SetDeadline(time.Time{})
	p.mu.Lock()
	p.errorText = ""
	p.mu.Unlock()
	writerDone := make(chan struct{})
	go func() { defer close(writerDone); p.sourceCommands(conn) }()
	defer func() { conn.Close(); <-writerDone }()
	for {
		record, err := readPreviewRecord(conn)
		if err != nil {

			p.sourceError("Önizleme eklentisi bağlantısı kesildi veya geçersiz veri gönderdi")
			return
		}
		if err := p.consume(record); err != nil {
			p.sourceError("Geçersiz H.264 önizleme akışı")
			return
		}
	}
}
func (p *previewService) sourceError(message string) {
	p.mu.Lock()
	p.errorText = message
	p.mu.Unlock()
}
func (p *previewService) sourceCommands(conn *net.UnixConn) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	active := false
	activeSince := time.Time{}
	for {
		p.mu.Lock()
		want := p.activeLocked() > 0
		p.mu.Unlock()
		command := byte('P')
		if want != active {
			if want {
				command = 'S'
				activeSince = time.Now()
				p.mu.Lock()
				p.lastKey = time.Time{}
				p.lastKeySent = time.Now()
				p.mu.Unlock()
			} else {
				command = 'T'
			}
			active = want
		}
		if active {
			p.mu.Lock()
			last := p.lastFrame
			p.mu.Unlock()
			if last.Before(activeSince) {
				last = activeSince
			}
			if time.Since(last) > 10*time.Second {
				p.sourceError("Eklentiden video alınamıyor")
				conn.Close()
				return
			}
		}
		if command == 'P' && want {
			p.mu.Lock()
			if !p.lastKey.IsZero() && time.Since(p.lastKeySent) >= time.Second {
				command = 'K'
				p.lastKey = time.Time{}
				p.lastKeySent = time.Now()
			}
			p.mu.Unlock()
		}
		conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
		if _, err := conn.Write([]byte{command}); err != nil {
			conn.Close()
			return
		}
		select {
		case <-p.ctx.Done():
			conn.Close()
			return
		case <-tick.C:
		case <-p.wake:
		}
	}
}
func drainSamples(ch chan media.Sample) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}
func (p *previewService) consume(record previewRecord) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if record.kind == 2 {
		config, err := parsePreviewConfig(record)
		if err != nil {
			return err
		}
		if p.generation != 0 && record.generation <= p.generation {
			return fmt.Errorf("non-increasing generation")
		}
		p.config = config
		p.generation = record.generation
		p.havePTS = false
		p.lastFrame = time.Time{}
		for _, v := range p.viewers {
			v.waiting = true
			drainSamples(v.frames)
		}
		return nil
	}
	if record.generation != p.generation || p.generation == 0 || !validPreviewFrame(record) {
		return fmt.Errorf("video without configuration")
	}
	// Wire prototype promises B-frame-free H.264. Reject reordering/overflow.
	if record.pts != record.dts || (p.havePTS && record.pts <= p.lastPTS) {
		return fmt.Errorf("invalid video timestamp")
	}
	if p.havePTS {
		delta := uint64(record.pts) - uint64(p.lastPTS)
		if delta > uint64(record.denominator)*10/uint64(record.numerator) {
			return fmt.Errorf("timestamp discontinuity")
		}
	}
	p.lastPTS = record.pts
	p.havePTS = true
	p.lastFrame = time.Now()
	duration := time.Duration(float64(time.Second) * float64(p.config.fpsDen) / float64(p.config.fpsNum))
	payload := record.data
	if record.key {
		payload = make([]byte, 0, len(p.config.headers)+len(record.data))
		payload = append(payload, p.config.headers...)
		payload = append(payload, record.data...)
	}
	for _, v := range p.viewers {
		if !v.ready || (v.waiting && !record.key) {
			continue
		}
		sample := media.Sample{Data: payload, Duration: duration}
		select {
		case v.frames <- sample:
			v.waiting = false
		default:
			drainSamples(v.frames)
			v.waiting = true // resume only at a later IDR
		}
	}
	return nil
}
func (p *previewService) removeViewer(id string) {
	p.mu.Lock()
	v := p.viewers[id]
	if v != nil {
		delete(p.viewers, id)
		close(v.done)
	}
	p.mu.Unlock()
	if v != nil {
		v.pc.Close()
		p.notify()
	}
}
func (p *previewService) requestKey() {
	p.mu.Lock()
	if p.lastKey.IsZero() || time.Since(p.lastKey) > time.Second {
		p.lastKey = time.Now()
	}
	p.mu.Unlock()
	p.notify()
}
func (p *previewService) routes(router *gin.Engine) {
	router.GET("/api/preview/status", func(c *gin.Context) { c.Header("Cache-Control", "no-store"); c.JSON(200, p.status()) })
	router.DELETE("/api/preview/session/:id", func(c *gin.Context) {
		if !pairingSameOrigin(c.Request) {
			c.Status(403)
			return
		}
		p.removeViewer(c.Param("id"))
		c.Status(204)
	})
	router.POST("/api/preview/offer", p.offer)
}
func (p *previewService) offer(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if !pairingSameOrigin(c.Request) {
		c.Status(403)
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)
	var offer webrtc.SessionDescription
	decoder := json.NewDecoder(c.Request.Body)
	if decoder.Decode(&offer) != nil || offer.Type != webrtc.SDPTypeOffer || len(offer.SDP) > 60*1024 {
		c.Status(400)
		return
	}
	if decoder.Decode(new(any)) != io.EOF {
		c.Status(400)
		return
	}
	p.mu.Lock()
	if p.ctx.Err() != nil || len(p.viewers) >= previewMaxViewers {
		p.mu.Unlock()
		c.Status(503)
		return
	}
	// Hold reservation while creating peer; API construction performs no network IO.
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		p.mu.Unlock()
		c.Status(503)
		return
	}
	engine := &webrtc.MediaEngine{}
	err := engine.RegisterCodec(webrtc.RTPCodecParameters{RTPCodecCapability: webrtc.RTPCodecCapability{
		MimeType: webrtc.MimeTypeH264, ClockRate: 90000, SDPFmtpLine: "level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f",
		RTCPFeedback: []webrtc.RTCPFeedback{{Type: "nack"}, {Type: "nack", Parameter: "pli"}},
	}, PayloadType: 102}, webrtc.RTPCodecTypeVideo)
	if err != nil {
		p.mu.Unlock()
		c.Status(500)
		return
	}
	// Host candidates only: no external STUN/TURN or internet dependency.
	settings := webrtc.SettingEngine{}
	settings.SetIncludeLoopbackCandidate(true)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(engine), webrtc.WithSettingEngine(settings))
	pc, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		p.mu.Unlock()
		c.Status(500)
		return
	}
	track, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000}, "program", "obs-preview")
	if err != nil {
		p.mu.Unlock()
		pc.Close()
		c.Status(500)
		return
	}
	sender, err := pc.AddTrack(track)
	if err != nil {
		p.mu.Unlock()
		pc.Close()
		c.Status(500)
		return
	}
	v := &previewViewer{id: hex.EncodeToString(raw), pc: pc, track: track, frames: make(chan media.Sample, 2), done: make(chan struct{}), waiting: true}
	p.viewers[v.id] = v
	p.mu.Unlock()
	success := false
	defer func() {
		if !success {
			p.removeViewer(v.id)
		}
	}()
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateConnected {
			p.mu.Lock()
			if p.viewers[v.id] == v {
				v.ready = true
				v.waiting = true
			}
			p.mu.Unlock()
			p.requestKey()
		} else if state == webrtc.PeerConnectionStateFailed || state == webrtc.PeerConnectionStateDisconnected {
			go p.removeViewer(v.id)
		}
	})
	go func() {
		for {
			packets, _, err := sender.ReadRTCP()
			if err != nil {
				return
			}
			for _, packet := range packets {
				if _, ok := packet.(*rtcp.PictureLossIndication); ok {
					p.requestKey()
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-v.done:
				return
			case sample := <-v.frames:
				if v.track.WriteSample(sample) != nil {
					p.removeViewer(v.id)
					return
				}
			}
		}
	}()
	request := c.Request.Clone(context.Background())
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		deadline := time.Now().Add(15 * time.Second)
		expiry := time.Now().Add(time.Hour)
		for {
			select {
			case <-v.done:
				return
			case <-p.ctx.Done():
				p.removeViewer(v.id)
				return
			case <-ticker.C:
				p.mu.Lock()
				ready := v.ready
				p.mu.Unlock()
				if (!ready && time.Now().After(deadline)) || time.Now().After(expiry) || (p.authorize != nil && !p.authorize(request)) {
					p.removeViewer(v.id)
					return
				}
			}
		}
	}()
	if pc.SetRemoteDescription(offer) != nil {
		c.Status(400)
		return
	}
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		c.Status(400)
		return
	}
	gathered := webrtc.GatheringCompletePromise(pc)
	if pc.SetLocalDescription(answer) != nil {
		c.Status(500)
		return
	}
	select {
	case <-gathered:
	case <-time.After(5 * time.Second):
		c.Status(504)
		return
	case <-c.Request.Context().Done():
		return
	case <-p.ctx.Done():
		c.Status(503)
		return
	}
	local := pc.LocalDescription()
	if local == nil {
		c.Status(500)
		return
	}
	c.JSON(200, struct {
		Type string `json:"type"`
		SDP  string `json:"sdp"`
		ID   string `json:"id"`
	}{local.Type.String(), local.SDP, v.id})
	success = true
}
