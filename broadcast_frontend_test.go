package main

import (
	"net/http/httptest"
	lib "obs-lan-controller/Lib"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestBroadcastStructure(t *testing.T) {
	h, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/", "/deck", "/yayin-yonetimi"} {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest("GET", path, nil))
		body := rr.Body.String()
		if strings.Contains(body, "broadcast-confirm") || strings.Contains(body, `id="confirm-dialog"`) {
			t.Errorf("confirmation UI on %s", path)
		}
		for _, marker := range []string{`id="broadcast-dashboard"`, `data-broadcast-action="start-stream"`, `data-broadcast-action="stop-record"`, `id="broadcast-log"`} {
			if strings.Contains(body, marker) != (path == "/yayin-yonetimi") {
				t.Errorf("%s marker %s", path, marker)
			}
		}
	}
}

func TestBroadcastFrontendBehavior(t *testing.T) {
	app, err := os.ReadFile("web/static/broadcast.js")
	if err != nil {
		t.Fatal(err)
	}
	script := `const assert=require('node:assert/strict');
const els=new Map();const e=id=>{if(!els.has(id))els.set(id,{textContent:'',disabled:true,hidden:false,dataset:{},classList:{toggle(){}},setAttribute(){},addEventListener(k,f){this[k]=f},appendChild(c){this.children.unshift(c)},prepend(c){this.children.unshift(c)},children:[],remove(){},focus(){},showModal(){this.open=true},close(){this.open=false}});return els.get(id)};
const actions=['start-stream','stop-stream','start-record','stop-record','pause-record','resume-record'].map(action=>Object.assign(e(action),{dataset:{broadcastAction:action}}));
global.document={querySelector:s=>e(s.slice(1)),querySelectorAll:()=>actions,createElement:()=>({remove(){const list=e('broadcast-log').children;list.splice(list.indexOf(this),1)}})};global.window={};let requests=[];let finish;global.fetch=(url,options)=>new Promise(resolve=>{requests.push({url,options});finish=resolve});
` + string(app) + `
(async()=>{
const status={connected:true,streamActive:false,recordAvailable:true,recordActive:false,fps:60,cpu:5,renderTime:1,renderSkipped:0,encodeSkipped:0};
window.OBSBroadcast.render(status);
assert(!e('start-stream').disabled);assert(e('stop-stream').disabled);assert(!e('start-record').disabled);
const pending=e('start-stream').click();assert.equal(requests.length,1);assert(actions.every(b=>b.disabled));
assert.deepEqual(JSON.parse(requests[0].options.body),{action:'start-stream',confirmed:true});
await e('start-stream').click();assert.equal(requests.length,1);
finish({ok:true,json:async()=>({...status,streamActive:true,duration:4000})});await pending;
assert(e('start-stream').disabled);assert(!e('stop-stream').disabled);assert.equal(e('manage-stream-state').textContent,'CANLI');
window.OBSBroadcast.render(null);assert(actions.every(b=>b.disabled));await e('stop-stream').click();assert.equal(requests.length,1);
window.OBSBroadcast.render({...status,recordAvailable:false});assert(e('start-record').disabled);
window.OBSBroadcast.render({...status,recordActive:true,recordPauseSupported:false});assert(e('pause-record').disabled);assert(e('resume-record').disabled);assert(e('pause-record').title.includes('desteklenmiyor'));assert(e('resume-record').title.includes('Yayınla aynı'));assert(e('record-pause-note').textContent.includes('paylaşılan encoder'));
window.OBSBroadcast.render({...status,recordActive:true,recordPaused:true});assert(!e('resume-record').disabled);assert(e('pause-record').disabled);
const failed=e('resume-record').click();finish({ok:false,status:409});await failed;assert(actions.every(b=>b.disabled));
for (const action of ['start-stream','stop-stream','start-record','stop-record','pause-record','resume-record']) {
 window.OBSBroadcast.render({...status,streamActive:action==='stop-stream',recordActive:['stop-record','pause-record','resume-record'].includes(action),recordPaused:action==='resume-record'});
 const before=requests.length;const operation=e(action).click();assert.equal(requests.length,before+1);assert.equal(JSON.parse(requests.at(-1).options.body).action,action);finish({ok:true,json:async()=>status});await operation;
}
console.log('direct click for all six actions, lock, verified state, disconnect and conflict OK');
})().catch(err=>{console.error(err);process.exitCode=1});`
	out, err := exec.Command("node", "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}
