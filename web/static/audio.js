"use strict";
(() => {
  const list = document.querySelector("#audio-inputs");
  const feedback = document.querySelector("#audio-feedback");
  if (!list || !feedback) return;
  let online = false, generation = 0, polling = false, revision = 0;
  const rows = new Map();
  const disable = () => rows.forEach(row => { row.slider.disabled = true; row.mute.disabled = true; });
  function connection(value) {
    if (online === value) return;
    online = value;
    if (!online) {
      generation++; revision++;
      rows.forEach(row => { clearTimeout(row.timer); row.timer = null; row.queued = null; row.dragging = false; row.busy = false; });
      disable(); feedback.textContent = "OBS bağlantısı yok";
    } else { feedback.textContent = "Ses kaynakları okunuyor…"; poll(); }
  }
  function update(row, input) {
    if (!row.dragging && !row.busy && row.queued === null && !row.timer) row.slider.value = String(Math.round(input.volume * 100));
    if (!row.busy) { row.muted = input.muted; row.mute.textContent = input.muted ? "Sesi aç" : "Sustur"; row.mute.setAttribute("aria-pressed", String(input.muted)); }
    row.slider.disabled = !online; row.mute.disabled = !online || row.busy;
  }
  function make(input) {
    const el = document.createElement("div"); el.className = "audio-row";
    const label = document.createElement("label"); label.textContent = input.inputName;
    const slider = document.createElement("input"); slider.type = "range"; slider.min = "0"; slider.max = "100"; slider.step = "1"; slider.setAttribute("aria-label", input.inputName + " ses seviyesi");
    const mute = document.createElement("button"); mute.type = "button"; mute.className = "btn btn-outline-warning"; mute.setAttribute("aria-label", input.inputName + " susturma");
    el.append(label, slider, mute); list.append(el);
    const row = {el, slider, mute, name: input.inputName, dragging:false, busy:false, queued:null, timer:null, muted:false};
    slider.addEventListener("pointerdown", () => { row.dragging = true; });
    slider.addEventListener("input", () => { if (!online) return; row.queued = Number(slider.value) / 100; revision++; schedule(row); });
    const release = () => { if (!row.dragging) return; row.dragging = false; if (!online) return; row.queued = Number(slider.value) / 100; revision++; clearTimeout(row.timer); row.timer = null; send(row); };
    slider.addEventListener("pointerup", release); slider.addEventListener("pointercancel", release); slider.addEventListener("change", release);
    mute.addEventListener("click", () => { if (online && !row.busy) send(row, !row.muted); });
    rows.set(input.inputName, row); return row;
  }
  function schedule(row) { if (!row.timer && !row.busy) row.timer = setTimeout(() => { row.timer = null; send(row); }, 150); }
  async function send(row, muteValue) {
    if (!online || row.busy) return;
    const isMute = typeof muteValue === "boolean";
    if (!isMute && row.queued === null) return;
    const value = row.queued; if (!isMute) row.queued = null;
    row.busy = true; revision++; const epoch = generation;
    row.mute.disabled = true;
    try {
      const payload = isMute ? {action:"mute", inputName:row.name, muted:muteValue} : {action:"volume", inputName:row.name, volume:value};
      const response = await fetch("/api/audio/action", {method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify(payload)});
      if (!response.ok) throw new Error();
      const actual = await response.json();
      if (epoch !== generation || !online) return;
      if (!actual || actual.inputName !== row.name || !Number.isFinite(actual.volume) || typeof actual.muted !== "boolean" || (isMute ? actual.muted !== muteValue : Math.abs(actual.volume - value) > 0.0001)) throw new Error();
      row.busy = false; update(row, actual); feedback.textContent = "Ses durumu doğrulandı";
    } catch (_) { if (epoch === generation && online) { feedback.textContent = "Ses işlemi doğrulanamadı"; row.busy = false; } }
    finally {
      if (epoch === generation && online) { row.busy = false; row.mute.disabled = false; if (row.queued !== null) schedule(row); }
    }
  }
  async function poll() {
    if (!online || polling) return;
    polling = true; const epoch = generation, started = revision;
    try {
      const response = await fetch("/api/audio", {cache:"no-store"}); if (!response.ok) throw new Error();
      const data = await response.json();
      if (!online || epoch !== generation || started !== revision) return;
      if (!Array.isArray(data.inputs)) throw new Error();
      const names = new Set();
      data.inputs.forEach(input => {
        if (typeof input.inputName !== "string" || !Number.isFinite(input.volume) || typeof input.muted !== "boolean") return;
        names.add(input.inputName); update(rows.get(input.inputName) || make(input), input);
      });
      rows.forEach((row, name) => { if (!names.has(name) && !row.busy && !row.dragging && row.queued === null) { row.el.remove(); rows.delete(name); } });
      feedback.textContent = names.size ? "" : "Ses kaynağı bulunamadı";
    } catch (_) { if (online && epoch === generation) { disable(); feedback.textContent = "Ses kaynakları okunamadı"; } }
    finally { polling = false; }
  }
  window.OBSAudio = {connection};
  setInterval(poll, 1000);
})();
