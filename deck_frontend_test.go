package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/http/httptest"
	lib "obs-lan-controller/Lib"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDeckStructure(t *testing.T) {
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest("GET", "/deck", nil))
	for _, want := range []string{`id="deck-dashboard"`, `id="deck-grid"`, `id="deck-edit"`, `id="deck-editor"`, `id="deck-feedback"`, `/static/deck.js`, `/static/deck.css`} {
		if !strings.Contains(rr.Body.String(), want) {
			t.Errorf("Deck missing %s", want)
		}
	}
	for _, forbidden := range []string{"data-broadcast-action", "Yayını başlat", "Kaydı başlat", "Özelleştirilebilir tuş ızgarası bu alanda"} {
		if strings.Contains(rr.Body.String(), forbidden) {
			t.Errorf("Deck contains %s", forbidden)
		}
	}
}

func TestDeckModel(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Fatal(err)
	}
	script := `const assert=require('node:assert/strict');global.window={};global.document={querySelector:()=>null};require('./web/static/deck.js');const m=window.OBSDeckModel;
 const state={connected:true,currentScene:'Main',scenes:['Main','Break'],items:[{sceneName:'Main',sceneItemId:0,sourceName:'Logo',targetId:'identity-logo',enabled:true}],inputs:[{inputName:'Mic',muted:false},{inputName:'Clip',mediaState:'OBS_MEDIA_STATE_PLAYING'}]};
 const slots=m.defaults(state);assert.equal(slots.length,12);assert.equal(slots.filter(Boolean).length,5);assert.equal(slots[0].target,'Main');
 assert.equal(slots[3].targetId,'identity-logo');
 const visibility={action:'visibility',sceneName:'Main',sceneItemId:0,targetId:'identity-logo'};
 assert.equal(m.active(slots[0],state),true);assert.deepEqual(m.payload(visibility,state),{...visibility,enabled:false});
 assert.equal(m.payload({...visibility,targetId:undefined},state),null,'legacy key must require explicit rebind');
 assert.equal(m.active({...visibility,targetId:undefined},state),false);
 for(const targetId of ['identity-other','',null])assert.equal(m.payload({...visibility,targetId},state),null,'stale collection/source identity');
 assert.equal(m.payload(visibility,{...state,items:[{...state.items[0],targetId:'identity-replacement'}]}),null);
 assert.deepEqual(m.payload({action:'mute',target:'Mic'},state),{action:'mute',inputName:'Mic',muted:true});assert.equal(m.payload({action:'mute',target:'Missing'},state),null);
 assert.equal(m.payload({action:'start-stream',target:'Main'},state),null);assert.equal(m.valid({action:'start-record'}),false);assert.equal(m.valid({action:'scene',target:'Main',color:'amber',icon:'scene',label:'Main'}),true);
 assert.equal(m.payload(slots[0],{...state,connected:false}),null);assert.equal(m.active(slots[0],{...state,connected:false}),false);
 assert.equal(m.valid({action:'scene',target:'Main',color:'evil',icon:'scene',label:'Main'}),false); console.log('defaults, allowlist, explicit state commands, missing targets, disconnect OK');`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	t.Log(string(out))
}

// Native selects and real CSS geometry catch failures a DOM double cannot model.
func TestDeckBrowserRegression(t *testing.T) {
	chrome := os.Getenv("CHROME_BIN")
	if chrome == "" {
		chrome = "/root/.hermes/tools/chromium-1208/chrome-linux64/chrome"
	}
	if _, err := os.Stat(chrome); err != nil {
		t.Skip("set CHROME_BIN to run real browser regressions")
	}
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/deck" {
			h.ServeHTTP(w, r)
			return
		}
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, r)
		body := strings.Replace(rr.Body.String(), "<head>", "<head><script>"+deckBrowserFixture+"</script>", 1)
		body = strings.Replace(body, "</body>", "<script>"+deckBrowserChecks+"</script></body>", 1)
		fmt.Fprint(w, body)
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	profile := t.TempDir()
	cmd := exec.CommandContext(ctx, chrome, "--headless", "--no-sandbox", "--disable-gpu", "--disable-background-networking", "--remote-debugging-port=0", "--user-data-dir="+profile, "about:blank")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
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
		t.Fatal("Chromium debugging port unavailable")
	}
	response, err := http.Get("http://127.0.0.1:" + port + "/json/list")
	if err != nil {
		t.Fatal(err)
	}
	var tabs []struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
		Type                 string `json:"type"`
		URL                  string `json:"url"`
	}
	json.NewDecoder(response.Body).Decode(&tabs)
	response.Body.Close()
	var debuggerURL string
	for _, tab := range tabs {
		if tab.Type == "page" && tab.URL == "about:blank" {
			debuggerURL = tab.WebSocketDebuggerURL
			break
		}
	}
	if debuggerURL == "" {
		t.Fatal("no Chromium about:blank page target")
	}
	socket, _, err := websocket.DefaultDialer.Dial(debuggerURL, nil)
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
	for _, viewport := range []struct {
		width, height int
		query         string
	}{{390, 844, ""}, {390, 844, "?legacy"}, {1280, 720, ""}, {844, 390, ""}} {
		call("Emulation.setDeviceMetricsOverride", map[string]any{"width": viewport.width, "height": viewport.height, "deviceScaleFactor": 1, "mobile": false})
		call("Runtime.evaluate", map[string]any{"expression": "document.querySelector('#regression-results')?.remove()"})
		call("Page.navigate", map[string]any{"url": server.URL + "/deck" + viewport.query})
		var result string
		for i := 0; i < 100; i++ {
			raw := call("Runtime.evaluate", map[string]any{"expression": "document.querySelector('#regression-results')?.textContent || ''", "returnByValue": true})
			var value struct {
				Result struct {
					Value string `json:"value"`
				} `json:"result"`
			}
			json.Unmarshal(raw, &value)
			result = value.Result.Value
			if result != "" {
				break
			}
			time.Sleep(30 * time.Millisecond)
		}
		if result != "PASS" {
			diagnostic := call("Runtime.evaluate", map[string]any{"expression": "JSON.stringify({url:location.href,ready:document.readyState,result:document.querySelector('#regression-results')?.outerHTML,fixture:typeof fixture,nextPoll:typeof nextPoll,deck:!!document.querySelector('#deck-grid')})", "returnByValue": true})
			t.Fatalf("browser regression %dx%d%s: %s; diagnostic: %s", viewport.width, viewport.height, viewport.query, result, diagnostic)
		}
		t.Logf("Chromium %dx%d%s: PASS", viewport.width, viewport.height, viewport.query)
	}
}

const deckBrowserFixture = `
window.fixture={connected:true,currentScene:'Main',scenes:['Main','Break'],items:[{sceneName:'Main',sceneItemId:0,sourceName:'Logo',targetId:'identity-logo',enabled:true}],inputs:[{inputName:'Mic',muted:false},{inputName:'Clip',mediaState:'OBS_MEDIA_STATE_PLAYING'}]};
window.requests=[];window.setTimeout=fn=>{window.nextPoll=fn;};window.fetch=async(url,options)=>{requests.push({url,options});return {ok:true,json:async()=>structuredClone(fixture)}};
localStorage.clear();
if(location.search.includes('legacy')) {
 const keys=Array(12).fill(null);
 keys[0]={action:'scene',target:'Main',color:'amber',icon:'scene',label:''};
 keys[1]={...keys[0],target:'Break'};
 keys[2]={action:'mute',target:'Mic',color:'green',icon:'mic',label:''};
 keys[3]={action:'visibility',sceneName:'Main',sceneItemId:0,color:'blue',icon:'eye',label:'Legacy logo'};
 keys[4]={action:'media-restart',target:'Clip',color:'purple',icon:'restart',label:''};
 localStorage.setItem('obs-controller.deck.v1',JSON.stringify(keys));
}
`

const deckBrowserChecks = `
document.addEventListener('DOMContentLoaded',async()=>{
 const result=document.createElement('pre');result.id='regression-results';document.body.append(result);
 const check=(ok,message)=>{if(!ok)throw Error(message)};
 try {
  for(let i=0;i<20;i++)await Promise.resolve();
  const errors=[];
  const verify=async(name,fn)=>{try{await fn()}catch(e){errors.push(name+': '+e.message)}};
  await verify('portrait',async()=>{
  const box=s=>document.querySelector(s).getBoundingClientRect();
  const grid=box('#deck-grid'),last=box('#deck-grid button:last-child'),bottom=box('.deck-bottom'),workspace=box('.deck-workspace'),sidebar=box('.deck-sidebar');
  check(last.bottom<=grid.bottom+1,'last key overflows grid: '+JSON.stringify({grid:grid.toJSON(),last:last.toJSON()}));
  check(grid.bottom<=bottom.top+1,'grid overlaps bottom controls');
  check(bottom.bottom<=workspace.bottom+1,'bottom controls overflow workspace');
  if(matchMedia('(max-width:700px) and (orientation:portrait)').matches)check(workspace.bottom<=sidebar.top+1,'workspace overlaps sidebar');
  else check(workspace.right<=sidebar.left+1,'workspace overlaps sidebar horizontally');
  check(document.documentElement.scrollWidth<=innerWidth,'horizontal document overflow');
  });
  const $=id=>document.getElementById(id);
  const settle=async()=>{for(let i=0;i<20;i++)await Promise.resolve()};
  if(location.search.includes('legacy')) {
   check(document.querySelectorAll('#deck-grid button')[3].disabled,'legacy key actionable before explicit rebind');
  }
  $('deck-edit').click();
  await verify('missing scene editor',async()=>{
   fixture.scenes=['Main'];await nextPoll();await settle();
   document.querySelectorAll('#deck-grid button')[1].click();
   check($('deck-target').selectedIndex===-1,'silently selected first remaining scene');
   check($('deck-save').disabled,'save enabled without explicit target');
   check($('deck-editor-feedback').textContent.length>0,'missing warning');
   $('deck-label').value='Label only';$('deck-form').dispatchEvent(new Event('submit',{cancelable:true}));
   check(JSON.parse(localStorage.getItem('obs-controller.deck.v1'))[1].target==='Break','label edit rebound missing scene');
   $('deck-target').selectedIndex=0;$('deck-target').dispatchEvent(new Event('change'));
   check(!$('deck-save').disabled,'explicit valid selection cannot recover');
   $('deck-form').dispatchEvent(new Event('submit',{cancelable:true}));
   check(JSON.parse(localStorage.getItem('obs-controller.deck.v1'))[1].target==='Main','explicit target not saved');
  });
  if($('deck-editor').open)$('deck-editor').close();
  await verify('stale visibility identity editor',async()=>{
   fixture.items[0].targetId='identity-replacement';await nextPoll();await settle();
   document.querySelectorAll('#deck-grid button')[3].click();
   check($('deck-target').selectedIndex===-1,'rebound reused scene item id');
   check($('deck-save').disabled,'save enabled for stale identity');
   $('deck-target').selectedIndex=0;$('deck-target').dispatchEvent(new Event('change'));
   $('deck-form').dispatchEvent(new Event('submit',{cancelable:true}));
   check(JSON.parse(localStorage.getItem('obs-controller.deck.v1'))[3].targetId==='identity-replacement','explicit identity not saved');
  });
  check(requests.every(r=>!r.options?.method),'editor issued an OBS action');
  if(errors.length)throw Error(errors.join(' | '));
  result.textContent='PASS';
 }catch(error){result.textContent='FAIL '+error.message}
});`
