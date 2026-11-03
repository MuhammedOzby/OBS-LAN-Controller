package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestAudioFrontendBehavior(t *testing.T) {
	app, err := os.ReadFile("web/static/audio.js")
	if err != nil {
		t.Fatal(err)
	}
	script := `const assert=require('node:assert/strict');
class Element { constructor(){this.children=[];this.events={};this.disabled=false;this.value='';this.textContent='';} append(...v){this.children.push(...v)} setAttribute(k,v){this[k]=v} addEventListener(k,v){this.events[k]=v} remove(){this.removed=true} }
const list=new Element(),feedback=new Element(); global.document={querySelector:s=>s==='#audio-inputs'?list:feedback,createElement:()=>new Element()};global.window={};
let timerId=0, timers=new Map(), interval; global.setTimeout=(f,ms)=>{let id=++timerId;timers.set(id,f);return id};global.clearTimeout=id=>timers.delete(id);global.setInterval=f=>{interval=f};
let requests=[];global.fetch=(url,options)=>new Promise((resolve,reject)=>requests.push({url,options,resolve,reject}));
const settle=async()=>{for(let i=0;i<8;i++)await Promise.resolve()};
const answer=(r,value)=>r.resolve({ok:true,json:async()=>value});
const source=(volume,muted=false)=>({inputs:[{inputName:'Mic',volume,muted}]});
` + string(app) + `
(async()=>{
window.OBSAudio.connection(true); assert.equal(requests.length,1); answer(requests[0],source(.5)); await settle();
const [label,slider,mute]=list.children[0].children; assert.equal(label.textContent,'Mic');assert.equal(slider.value,'50');assert.equal(slider.disabled,false);
interval(); answer(requests[1],source(.8,true));await settle();assert.equal(slider.value,'80');assert.equal(mute.textContent,'Sesi aç');
slider.events.pointerdown();slider.value='30';slider.events.input();slider.value='40';slider.events.input();assert.equal(timers.size,1);
const fire=()=>{const [id,fn]=timers.entries().next().value;timers.delete(id);fn()};fire();assert.equal(JSON.parse(requests[2].options.body).volume,.4);
slider.value='60';slider.events.input();slider.events.pointerup();assert.equal(requests.length,3);
interval();answer(requests[3],source(.4));await settle();assert.equal(slider.value,'60');
answer(requests[2],{inputName:'Mic',volume:.4,muted:true});await settle();assert.equal(timers.size,1);fire();assert.equal(JSON.parse(requests[4].options.body).volume,.6);
window.OBSAudio.connection(false);assert(slider.disabled);assert(mute.disabled);assert.equal(feedback.textContent,'OBS bağlantısı yok');
window.OBSAudio.connection(true);answer(requests[5],source(.2,false));await settle();assert.equal(slider.value,'20');assert.equal(slider.disabled,false);
mute.events.click();assert.equal(requests.length,7);assert(mute.disabled);
answer(requests[4],{inputName:'Mic',volume:.6,muted:true});await settle();assert(mute.disabled);assert.equal(slider.value,'20');
 answer(requests[6],{inputName:'Mic',volume:.2,muted:true});await settle();assert.equal(mute.textContent,'Sesi aç');
 mute.events.click(); answer(requests[7],{inputName:'Mic',volume:.2}); await settle(); assert.equal(feedback.textContent,'Ses işlemi doğrulanamadı'); assert.equal(mute.textContent,'Sesi aç');
console.log('audio polling, external changes, coalescing, final value, stale responses and offline disable OK');
})().catch(e=>{console.error(e);process.exitCode=1});`
	out, err := exec.Command("node", "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}
