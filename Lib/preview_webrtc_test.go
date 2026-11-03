//go:build linux

package lib

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
)

// Real ICE/DTLS/SRTP transport, simulated Annex-B source. Not a decoder test.
func TestPreviewWebRTCRoundTrip(t *testing.T) {
	runPreviewRoundTrip(t, false)
}

func TestPreviewEncodedVideoRoundTrip(t *testing.T) {
	if os.Getenv("PREVIEW_FFMPEG_TEST") != "1" {
		t.Skip("explicit ffmpeg encoder/decoder integration")
	}
	runPreviewRoundTrip(t, true)
}

func runPreviewRoundTrip(t *testing.T, encoded bool) {
	var headers, frame []byte
	if encoded {
		command := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-f", "lavfi", "-i", "testsrc2=size=1280x720:rate=30", "-frames:v", "1", "-c:v", "libx264", "-profile:v", "baseline", "-preset", "veryfast", "-tune", "zerolatency", "-f", "h264", "pipe:1")
		var stderr bytes.Buffer
		command.Stderr = &stderr
		data, err := command.Output()
		if err != nil {
			t.Fatalf("encode: %v %s", err, stderr.String())
		}
		for _, nal := range annexBNALs(data) {
			start := []byte{0, 0, 0, 1}
			if nal[0]&31 == 7 || nal[0]&31 == 8 {
				headers = append(headers, start...)
				headers = append(headers, nal...)
			} else {
				frame = append(frame, start...)
				frame = append(frame, nal...)
			}
		}
	} else {
		frame = []byte{0, 0, 0, 1, 0x65, 0x88, 0x84}
	}
	p, path := setupPreview(t)
	source := connectTestSource(t, path)
	handler, err := newHandlerWithPreview(nil, p)
	if err != nil {
		t.Fatal(err)
	}
	browser, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	defer browser.Close()
	if _, err := browser.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly}); err != nil {
		t.Fatal(err)
	}
	received := make(chan []byte, 1)
	browser.OnTrack(func(track *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		depay := &codecs.H264Packet{}
		var video []byte
		for {
			packet, _, err := track.ReadRTP()
			if err != nil {
				return
			}
			nal, err := depay.Unmarshal(packet.Payload)
			if err != nil {
				return
			}
			video = append(video, nal...)
			if packet.Marker {
				select {
				case received <- video:
				default:
				}
				return
			}
		}
	})
	offer, err := browser.CreateOffer(nil)
	if err != nil {
		t.Fatal(err)
	}
	gathered := webrtc.GatheringCompletePromise(browser)
	if err := browser.SetLocalDescription(offer); err != nil {
		t.Fatal(err)
	}
	select {
	case <-gathered:
	case <-time.After(5 * time.Second):
		t.Fatal("offer gather timeout")
	}
	body, _ := json.Marshal(browser.LocalDescription())
	req := httptest.NewRequest("POST", "http://example.test/api/preview/offer", bytes.NewReader(body))
	req.Header.Set("Origin", "http://example.test")
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("offer: %d %s", response.Code, response.Body)
	}
	var answer struct {
		Type string `json:"type"`
		SDP  string `json:"sdp"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &answer); err != nil {
		t.Fatal(err)
	}
	if err := browser.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer.SDP}); err != nil {
		t.Fatal(err)
	}
	awaitCommand(t, source, 'S')
	config := testConfig()
	if encoded {
		config = append(config[:16:16], headers...)
	}
	if _, err := source.Write(testRecord(2, 1, false, 0, config)); err != nil {
		t.Fatal(err)
	}
	// Packetizer may hold one access unit; send several independent IDRs.
	for i := int64(0); i < 3; i++ {
		if _, err := source.Write(testRecord(3, 1, true, i, frame)); err != nil {
			t.Fatal(err)
		}
		time.Sleep(40 * time.Millisecond)
	}
	select {
	case payload := <-received:
		if len(payload) == 0 {
			t.Fatal("empty RTP")
		}
		if encoded {
			output := filepath.Join(t.TempDir(), "decoded.yuv")
			command := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error", "-f", "h264", "-i", "pipe:0", "-frames:v", "1", "-pix_fmt", "yuv420p", "-f", "rawvideo", output)
			command.Stdin = bytes.NewReader(payload)
			if log, err := command.CombinedOutput(); err != nil {
				t.Fatalf("decode: %v %s", err, log)
			}
			info, err := os.Stat(output)
			if err != nil || info.Size() != 1280*720*3/2 {
				t.Fatalf("decoded dimensions: %v %v", info, err)
			}
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no RTP received")
	}
	status := p.status()
	if !status.Streaming || status.Width != 1280 || status.Viewers != 1 {
		t.Fatalf("status: %+v", status)
	}
	req = httptest.NewRequest("DELETE", "http://example.test/api/preview/session/"+answer.ID, nil)
	req.Header.Set("Origin", "http://example.test")
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	if response.Code != 204 {
		t.Fatal(response.Code)
	}
	awaitCommand(t, source, 'T')
	if p.status().Viewers != 0 {
		t.Fatal("viewer not released")
	}
}
