package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	lib "obs-lan-controller/Lib"
)

func TestFrontendControlStructure(t *testing.T) {
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
	body := rr.Body.String()
	for _, want := range []string{`id="control-dashboard"`, `id="program-video"`, `autoplay muted playsinline`, `id="preview-placeholder"`, `id="preview-status"`, `id="preview-fps"`, `/static/preview.js`, `id="stream-state"`, "Yayın yapılmıyor", `id="obs-fps"`, `id="obs-memory"`, `id="obs-render"`, `id="obs-skips"`, `id="stat-cpu"`, `id="stat-network"`, `id="stat-drop"`, `id="stat-quality"`, `id="scene-feedback"`, "CPU", "AĞ", "DROP", "Kalite"} {
		if !strings.Contains(body, want) {
			t.Errorf("control missing %q", want)
		}
	}
	if got := strings.Count(body, `data-scene-slot=`); got != 3 {
		t.Errorf("scene slots = %d, want exactly 3", got)
	}
	if strings.Contains(body, "Yayını başlat") || strings.Contains(body, "Yayını durdur") {
		t.Error("broadcast controls on control page")
	}
	if strings.Count(body, `class="nav-link`) != 4 || strings.Contains(body, `id="program-preview"`) {
		t.Error("exactly four tabs and no screenshot element expected")
	}
	for _, path := range []string{"/deck", "/yayin-yonetimi"} {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest("GET", path, nil))
		if strings.Contains(rr.Body.String(), `id="control-dashboard"`) {
			t.Errorf("dashboard leaked to %s", path)
		}
		if strings.Contains(rr.Body.String(), `/static/preview.js`) {
			t.Errorf("preview loaded outside control page: %s", path)
		}
	}
}

// Run only explicitly: serves real embedded templates without reading configuration or connecting OBS.
func TestFrontendBrowserServer(t *testing.T) {
	addr := os.Getenv("FRONTEND_BROWSER_ADDR")
	if addr == "" {
		t.Skip("explicit browser harness only")
	}
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(http.ListenAndServe(addr, h))
}

func TestFrontendPollingAndFailure(t *testing.T) {
	runFrontendJS(t, `(async()=>{
 await settle();
 assert.equal(requests.length,1); assert.equal(requests[0].url,'/api/status?preview=0');
 assert.equal(timers.length,1); assert.equal(timers[0].ms,1000);
 assert.equal(slots[0].textContent,'<b>Bir</b>'); assert.equal(slots[1].attributes['aria-pressed'],'true');
 assert.equal(slots[2].disabled,false); assert.equal(element('stream-state').textContent,'Canlı');
 assert.equal(element('stream-duration').textContent,'00:01:05');
 assert.equal(element('stat-network').textContent,'3200 kb/s');
 assert.equal(element('program-preview').src,'','control polling never renders screenshots');
 replies.push(new Error('offline')); await timers.shift().fn();await settle();
 assert.equal(element('stat-cpu').textContent,'—'); assert.equal(element('obs-fps').textContent,'—');
 assert.equal(element('stat-quality').textContent,'Bilinmiyor');
 assert.equal(element('program-preview').src,'');
 assert(slots.every(s=>s.disabled));
 replies.push({...sample,streamActive:false,bitrateKbps:null,scenes:['Tek'],preview:'',previewAvailable:false});
 await timers.shift().fn(); await settle();
 assert.equal(element('stream-state').textContent,'Yayın yapılmıyor');
 assert.equal(element('stat-network').textContent,'—');assert.equal(element('stat-quality').textContent,'Bilinmiyor');
 assert.equal(slots[0].disabled,false); assert.equal(slots[1].disabled,true);
 assert.equal(element('program-preview').src,'');
 let resolvePoll; global.fetch=()=>new Promise(resolve=>{resolvePoll=resolve});
 const pending=timers.shift().fn(); await settle();
 assert.equal(timers.length,0,'no next request scheduled before current request settles');
 resolvePoll({ok:true,json:async()=>({...sample,connected:false,reconnecting:true})});await pending;await settle();
 assert.equal(timers.length,1);assert.equal(element('connection-status').textContent,'● OBS yeniden bağlanıyor');
 assert.equal(element('stat-quality').textContent,'Bilinmiyor');assert(slots.every(s=>s.disabled));
 console.log('sequential poll, safe names, metrics, failure reset, missing slots OK');
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendSceneCommand(t *testing.T) {
	runFrontendJS(t, `(async()=>{
 await settle(); assert.equal(typeof slots[0].click,'function');
 let finish; const originalFetch=global.fetch;
 global.fetch=(url,options)=>url==='/api/scene'?new Promise(resolve=>{requests.push({url,options});finish=resolve}):originalFetch(url,options);
 const pending=slots[0].click(); await settle();
 assert(slots.every(s=>s.disabled)); assert.equal(requests[1].url,'/api/scene');
 assert.deepEqual(JSON.parse(requests[1].options.body),{sceneName:'<b>Bir</b>'});
 assert.equal(requests[1].options.method,'POST');
 await slots[1].click(); assert.equal(requests.length,2);
 finish({ok:true,status:204}); await pending; await settle();
 assert.equal(slots[0].attributes['aria-pressed'],'false');
 assert.equal(element('scene-feedback').textContent,'Sahne değişimi doğrulandı');
 global.fetch=async()=>({ok:false,status:403}); await slots[0].click();
 assert.equal(element('scene-feedback').textContent,'İstek reddedildi; sayfayı yenileyin');
 assert.equal(slots[0].attributes['aria-pressed'],'false');
 global.fetch=async()=>({ok:false,status:400}); await slots[0].click();
 assert.equal(element('scene-feedback').textContent,'Sahne bulunamadı; liste yenilenecek');
 assert.equal(slots[0].attributes['aria-pressed'],'false');
 global.fetch=async()=>({ok:false,status:503}); await slots[0].click();
 assert(slots.every(s=>s.disabled)); assert.equal(element('stat-cpu').textContent,'—');
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendFavoriteSlots(t *testing.T) {
	runFrontendJS(t, `(async()=>{
 await settle();
 replies.push({...sample, scenes:['<b>Bir</b>','Second','Fourth'], favorites:['Fourth','','Missing']});
 await timers.shift().fn(); await settle();
 assert.equal(slots[0].textContent,'Fourth'); assert.equal(slots[0].disabled,false);
 assert.equal(slots[1].disabled,true); assert.equal(slots[2].disabled,true);
 assert.equal(slots[2].textContent,'Missing (kullanılamıyor)');
 let sent; global.fetch=async(url, options)=>{sent=JSON.parse(options.body); return {status:204}};
 await slots[0].click(); assert.deepEqual(sent,{sceneName:'Fourth'});
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendOtherPageNoSnapshot(t *testing.T) {
	runFrontendJS(t, `(async()=>{await settle(); global.otherPage=true; requests=[];`+
		func() string {
			b, e := os.ReadFile("web/static/app.js")
			if e != nil {
				t.Fatal(e)
			}
			return string(b)
		}()+
		`await settle(); assert.equal(requests[0].url,'/api/status?preview=0');})().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendDOMSelectorsAndDeckGuard(t *testing.T) {
	app, err := os.ReadFile("web/static/app.js")
	if err != nil {
		t.Fatal(err)
	}
	runFrontendJS(t, `(async()=>{
  await settle();
  assert.equal(document.querySelector('#does-not-exist'),null);
  assert.equal(document.querySelector('#deck-dashboard'),null);
  global.deckPage=true; requests=[]; timers.length=0;
 `+string(app)+`
  await settle(); assert.equal(requests.length,0,'Deck must own its polling');
  assert.equal(timers.length,0,'shared app must not schedule a second Deck poll');
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func runFrontendJS(t *testing.T, checks string) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Fatal("Node required for frontend behavior tests")
	}
	app, err := os.ReadFile("web/static/app.js")
	if err != nil {
		t.Fatal(err)
	}
	script := frontendDOMHarness + "\n" + string(app) + "\n" + checks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("frontend behavior: %v\n%s", err, out)
	}
}

func runPreviewJS(t *testing.T, checks string) {
	t.Helper()
	app, err := os.ReadFile("web/static/preview.js")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "node", "-e", previewDOMHarness+"\n"+string(app)+"\n"+checks).CombinedOutput()
	if err != nil {
		t.Fatalf("preview behavior: %v\n%s", err, out)
	}
}

func TestFrontendPreviewLifecycle(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 await settle();
 assert.equal(peers.length,1); assert.deepEqual(peers[0].transceiver,['video',{direction:'recvonly'}]);
 const offer=requests.find(r=>r.url==='/api/preview/offer');
 assert(offer); assert.equal(offer.options.method,'POST');
 assert.equal(offer.options.credentials,'same-origin'); assert(offer.options.signal);
 assert.deepEqual(JSON.parse(offer.options.body),{type:'offer',sdp:'complete-offer'});
 assert.equal(peers[0].remoteDescription.sdp,'answer');
 assert.equal(video.muted,true); assert.equal(video.playsInline,true);
 peers[0].ontrack({track:{kind:'video'},streams:[]}); await settle();
 assert(video.srcObject); video.listeners.playing();
 assert.equal(video.hidden,false);assert.equal(elements['preview-placeholder'].hidden,true);
 const obsText='OBS controls remain independent'; elements['connection-status']={textContent:obsText};
 document.hidden=true; document.listeners.visibilitychange(); await settle();
 assert.equal(peers[0].closed,true);assert.equal(video.srcObject,null);assert.equal(video.hidden,true);
 const close=requests.find(r=>r.options?.method==='DELETE');
 assert.equal(close.url,'/api/preview/session/viewer-1');assert.equal(close.options.keepalive,true);
 assert.equal(elements['connection-status'].textContent,obsText);
 assert.equal(timers.size,0,'hidden page has no preview work scheduled');
 document.hidden=false;document.listeners.visibilitychange();await settle();assert.equal(peers.length,2);
 window.listeners.pagehide();await settle();assert.equal(peers[1].closed,true);assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewPageGuards(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 await settle();window.listeners.pagehide();await settle();requests.length=0;
 const script=require('node:fs').readFileSync('web/static/preview.js','utf8');
 const peerCount=peers.length;
 document.hidden=true;eval(script);await settle();
 assert.equal(peers.length,peerCount);assert.equal(requests.length,0);
 document.hidden=false;const originalQuery=document.querySelector;
 document.querySelector=s=>s==='#control-dashboard'?null:originalQuery(s);
 eval(script);await settle();assert.equal(peers.length,peerCount);assert.equal(requests.length,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewInvalidAnswerCleanup(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 fetchImpl=async()=>({ok:true,json:async()=>({type:'offer',sdp:'invalid',id:'bad-viewer'})});
 await settle();assert.equal(peers[0].closed,true);assert.equal(peers[0].remoteDescription,undefined);
 assert(requests.some(r=>r.url==='/api/preview/session/bad-viewer'&&r.options.method==='DELETE'));
 window.listeners.pagehide();assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewStaleAnswer(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 let resolveOffer;
 fetchImpl=(url,options)=>url==='/api/preview/offer'?new Promise(resolve=>{resolveOffer=resolve}):Promise.resolve({ok:true,json:async()=>({})});
 await settle();assert(resolveOffer);
 document.hidden=true;document.listeners.visibilitychange();await settle();
 resolveOffer({ok:true,json:async()=>({type:'answer',sdp:'late',id:'late/id'})});await settle();
 assert.equal(peers[0].remoteDescription,undefined);
 assert(requests.some(r=>r.url==='/api/preview/session/late%2Fid'&&r.options.method==='DELETE'));
 assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewBoundedRetry(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 fetchImpl=async()=>({ok:false,status:503});
 await settle();
 for(const delay of [1000,2000,4000,8000])await fireTimer(delay);
 assert.equal(peers.length,5);assert(peers.every(pc=>pc.closed));assert.equal(timers.size,0);
 assert.equal(elements['preview-retry'].hidden,false);
 fetchImpl=async(url)=>({ok:true,json:async()=>url==='/api/preview/offer'?{type:'answer',sdp:'answer',id:'recovered'}:{connected:true,streaming:true}});
 elements['preview-retry'].listeners.click();await settle();assert.equal(peers.length,6);
 window.listeners.pagehide();await settle();assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewICETimeout(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 peers[0].iceGatheringState='gathering';await settle();
 assert.equal(requests.length,0,'must wait for nontrickle ICE');
 await fireTimer(5000);assert.equal(peers[0].closed,true);assert.equal(requests.length,0);
 window.listeners.pagehide();assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewFetchDeadline(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 let signal;
 fetchImpl=(url,options)=>new Promise((resolve,reject)=>{signal=options.signal;signal.addEventListener('abort',()=>reject(new Error('aborted')))});
 await settle();assert(signal);assert.equal(signal.aborted,false);
 await fireTimer(5000);assert.equal(signal.aborted,true);assert.equal(peers[0].closed,true);
 window.listeners.pagehide();assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewStatsAndAuth(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 await settle();peers[0].ontrack({track:{kind:'video'},streams:[]});await settle();video.listeners.playing();
 await fireTimer(1000);assert.equal(elements['preview-fps'].textContent,'30.0');
 fetchImpl=async()=>({ok:false,status:401});
 await fireTimer(1000);assert.equal(peers[0].closed,true);assert.equal(timers.size,0);
 assert.equal(elements['preview-fps'].textContent,'—');assert.equal(elements['preview-retry'].hidden,true);
 assert.match(elements['preview-status'].textContent,/reddedildi/);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

func TestFrontendPreviewStartupGrace(t *testing.T) {
	runPreviewJS(t, `(async()=>{
 fetchImpl=async(url)=>({ok:true,json:async()=>url==='/api/preview/offer'?{type:'answer',sdp:'answer',id:'starting'}:{connected:true,streaming:false}});
 await settle();
 assert.equal(peers[0].closed,undefined,'first status may precede ICE and encoder startup');
 await fireTimer(1000);
 assert.equal(peers[0].closed,undefined,'startup grace must survive a pending source');
 await fireTimer(12000);
 assert.equal(peers[0].closed,true,'video startup remains bounded');
 window.listeners.pagehide();assert.equal(timers.size,0);
 })().catch(e=>{console.error(e);process.exitCode=1});`)
}

const previewDOMHarness = `
const assert=require('node:assert/strict');
const elements={};
for(const id of ['program-video','preview-placeholder','preview-status','preview-fps','preview-retry'])elements[id]={hidden:false,textContent:'',listeners:{},addEventListener(k,fn){this.listeners[k]=fn}};
const video=elements['program-video'];video.play=async()=>{};video.pause=()=>{};
global.document={hidden:false,listeners:{},querySelector:s=>s==='#control-dashboard'?{}:elements[s.slice(1)]||null,addEventListener(k,fn){this.listeners[k]=fn}};
global.window={listeners:{},addEventListener(k,fn){this.listeners[k]=fn}};
let nextTimer=1;const timers=new Map();global.setTimeout=(fn,ms)=>{const id=nextTimer++;timers.set(id,{fn,ms});return id};global.clearTimeout=id=>timers.delete(id);
const peers=[];
class FakePeer {
 constructor(){this.iceGatheringState='complete';this.connectionState='new';this.listeners={};peers.push(this)}
 addTransceiver(...args){this.transceiver=args}
 async createOffer(){return {type:'offer',sdp:'initial-offer'}}
 async setLocalDescription(){this.localDescription={type:'offer',sdp:'complete-offer'}}
 async setRemoteDescription(answer){this.remoteDescription=answer}
 addEventListener(k,fn){this.listeners[k]=fn} removeEventListener(k){delete this.listeners[k]}
 async getStats(){return new Map([['video',{type:'inbound-rtp',kind:'video',framesDecoded:30,framesPerSecond:30,timestamp:1000}]])}
 close(){this.closed=true}
}
global.RTCPeerConnection=FakePeer;global.MediaStream=class{constructor(tracks){this.tracks=tracks}getTracks(){return this.tracks}};
const requests=[];let fetchImpl=async(url,options)=>({ok:true,status:200,json:async()=>url==='/api/preview/offer'?{type:'answer',sdp:'answer',id:'viewer-'+peers.length}:{connected:true,streaming:true,width:1280,height:720,fps:30}});
global.fetch=(url,options)=>{requests.push({url,options});return fetchImpl(url,options)};
const settle=async()=>{for(let i=0;i<50;i++)await Promise.resolve()};
const fireTimer=async(ms)=>{const entry=[...timers].find(([id,t])=>t.ms===ms);assert(entry,'timer '+ms);timers.delete(entry[0]);entry[1].fn();await settle()};
`

const frontendDOMHarness = `
const assert = require('node:assert/strict');
const elements = new Map();
function element(id) { if (!elements.has(id)) elements.set(id, {id, textContent:'', hidden:false, disabled:false, src:'', dataset:{}, attributes:{}, classList:{toggle(){}}, setAttribute(k,v){this.attributes[k]=v}, removeAttribute(k){delete this.attributes[k];if(k==='src')this.src=''}, addEventListener(k,v){this[k]=v}}); return elements.get(id); }
const slots = [0,1,2].map(i=>{const e=element('slot'+i); e.dataset.sceneSlot=String(i);return e});
const controlIds=['control-dashboard','stream-state','stream-duration','current-scene','obs-fps','obs-memory','obs-render','obs-skips','stat-cpu','stat-network','stat-drop','stat-quality','program-preview','preview-placeholder','scene-feedback'];
controlIds.forEach(element);element('connection-status');
global.document = {querySelector:s=>s==='.connection-status'?elements.get('connection-status'):s==='#deck-dashboard'?(global.deckPage?element('deck-dashboard'):null):global.otherPage?null:elements.get(s.replace(/^#/,''))||null, querySelectorAll:s=>s==='[data-scene-slot]'&&!global.otherPage?slots:[]};
global.window = {location:{pathname:'/'}};
const timers = []; global.setTimeout=(fn,ms)=>{timers.push({fn,ms}); return timers.length};global.clearTimeout=()=>{};
let requests=[]; const sample={connected:true,streamActive:true,reconnecting:false,duration:65000,bitrateKbps:3200,cpu:12.5,fps:59.94,memory:700,renderTime:2.5,renderSkipped:2,renderTotal:1000,encodeSkipped:1,encodeTotal:1000,dropped:3,total:1000,congestion:0,scenes:['<b>Bir</b>','İki','Üç','Dört'],currentScene:'İki',preview:'data:image/png;base64,AA==',previewAvailable:true}; let replies=[sample];
global.fetch=async(url,options)=>{requests.push({url,options});const next=replies.shift();if(next instanceof Error)throw next;return {ok:true,status:200,json:async()=>next,...next}};
const settle=async()=>{for(let i=0;i<15;i++)await Promise.resolve()};
`
