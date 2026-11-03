package lib

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestPairingRealMobileLayoutChrome(t *testing.T) {
	chrome := os.Getenv("CHROME_BIN")
	if chrome == "" {
		t.Skip("CHROME_BIN required")
	}
	listener, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	router, err := newHandler(nil)
	if err != nil {
		t.Fatal(err)
	}
	pairing := newPairingHandler(router)
	pairing.listenAddr = listener.Addr()
	if len(pairingAddresses(pairing.listenAddr)) == 0 {
		t.Skip("no non-loopback interface")
	}
	server, err := newHTTPServer(listener.Addr().String(), pairing)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	go server.Serve(listener)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	profile := t.TempDir()
	cmd := exec.CommandContext(ctx, chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--remote-debugging-port=0", "--user-data-dir="+profile, "about:blank")
	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { cmd.Process.Kill(); cmd.Wait() }()
	var port string
	for i := 0; i < 100; i++ {
		data, _ := os.ReadFile(filepath.Join(profile, "DevToolsActivePort"))
		if len(data) > 0 {
			port = strings.Split(string(data), "\n")[0]
			break
		}
		time.Sleep(30 * time.Millisecond)
	}
	if port == "" {
		t.Fatal("Chrome debugging port unavailable")
	}
	response, err := http.Get("http://127.0.0.1:" + port + "/json/list")
	if err != nil {
		t.Fatal(err)
	}
	var tabs []struct {
		Type                 string `json:"type"`
		URL                  string `json:"url"`
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	err = json.NewDecoder(response.Body).Decode(&tabs)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	var debug string
	for _, tab := range tabs {
		if tab.Type == "page" && tab.URL == "about:blank" {
			debug = tab.WebSocketDebuggerURL
			break
		}
	}
	if debug == "" {
		t.Fatal("missing page target")
	}
	socket, _, err := websocket.DefaultDialer.Dial(debug, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer socket.Close()
	socket.SetReadDeadline(time.Now().Add(20 * time.Second))
	id := 0
	call := func(method string, params any) json.RawMessage {
		id++
		if err := socket.WriteJSON(map[string]any{"id": id, "method": method, "params": params}); err != nil {
			t.Fatal(err)
		}
		for {
			var msg struct {
				ID     int             `json:"id"`
				Result json.RawMessage `json:"result"`
				Error  any             `json:"error"`
			}
			if err := socket.ReadJSON(&msg); err != nil {
				t.Fatal(err)
			}
			if msg.ID == id {
				if msg.Error != nil {
					t.Fatalf("CDP: %v", msg.Error)
				}
				return msg.Result
			}
		}
	}
	base := "http://127.0.0.1:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	for _, viewport := range []struct{ w, h int }{{390, 844}, {844, 390}, {1280, 720}} {
		call("Emulation.setDeviceMetricsOverride", map[string]any{"width": viewport.w, "height": viewport.h, "deviceScaleFactor": 1, "mobile": false})
		for _, path := range []string{"/mobil-cihaz-ile-giris", "/pair"} {
			call("Page.navigate", map[string]any{"url": base + path})
			var outcome string
			for i := 0; i < 100; i++ {
				raw := call("Runtime.evaluate", map[string]any{"expression": "location.pathname !== " + strconv.Quote(path) + " ? 'WAIT' : " + `(()=>{if(document.readyState!=='complete'||!document.querySelector('.pair-qr-list img'))return 'WAIT';try{const list=document.querySelector('.pair-qr-list'),img=list.querySelector('img');if(!img.complete||!img.naturalWidth)return 'WAIT';if(!list.dataset.added){const copy=list.firstElementChild.cloneNode(true);copy.querySelector('p').textContent='[ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff]:65535';list.append(copy);list.dataset.added='yes'}const panel=document.querySelector('.pair-panel'),root=panel||document.scrollingElement;if(document.documentElement.scrollWidth>innerWidth||root.scrollWidth>root.clientWidth+1)throw Error('horizontal overflow');if(img.getBoundingClientRect().width>root.clientWidth)throw Error('QR overflow');if(panel&&(document.querySelectorAll('.app-nav .nav-link').length!==4||document.querySelector('[data-broadcast-action]')))throw Error('navigation');for(const el of [document.querySelector('a[href="/pair"]'),document.querySelector('button[type="submit"]')]){if(!el)throw Error('missing control');el.scrollIntoView({block:'center'});const r=el.getBoundingClientRect(),p=panel?panel.getBoundingClientRect():{top:0,bottom:innerHeight};if(r.top<Math.max(0,p.top)-1||r.bottom>Math.min(innerHeight,p.bottom)+1)throw Error('control clipped')}return 'PASS'}catch(e){return 'FAIL '+e.message}})()`, "returnByValue": true})
				var result struct {
					Result struct {
						Value string `json:"value"`
					} `json:"result"`
				}
				json.Unmarshal(raw, &result)
				outcome = result.Result.Value
				if outcome != "WAIT" && outcome != "" {
					break
				}
				time.Sleep(30 * time.Millisecond)
			}
			if outcome != "PASS" {
				t.Fatalf("%dx%d %s: %s", viewport.w, viewport.h, path, outcome)
			}
		}
	}
}
