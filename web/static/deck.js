"use strict";
(() => {
  const actions = ["scene", "visibility", "mute", "media-play", "media-pause", "media-restart", "media-stop"];
  const colors = ["amber", "blue", "green", "purple"];
  const icons = {scene: "▣", eye: "◉", mic: "♩", play: "▶", pause: "Ⅱ", restart: "↻", stop: "■"};
  const titles = {scene: "SAHNE", visibility: "KAYNAK", mute: "SES", "media-play": "OYNAT", "media-pause": "DURAKLAT", "media-restart": "YENİDEN", "media-stop": "DURDUR"};
  const item = (key, s) => typeof key.targetId === "string" && key.targetId.length > 0 ? s?.items?.find(i => i.sceneName === key.sceneName && i.sceneItemId === key.sceneItemId && i.targetId === key.targetId) : undefined;
  const input = (key, s) => s?.inputs?.find(i => i.inputName === key.target);
  function valid(key) {
    return !!key && actions.includes(key.action) && colors.includes(key.color) && Object.hasOwn(icons, key.icon) && typeof key.label === "string" && key.label.length <= 40 && (key.action === "visibility" ? typeof key.sceneName === "string" && Number.isInteger(key.sceneItemId) && key.sceneItemId >= 0 : typeof key.target === "string" && key.target.length > 0);
  }
  function payload(key, s) {
    if (!key || !s?.connected || !actions.includes(key.action)) return null;
    if (key.action === "scene") return s.scenes?.includes(key.target) ? {action: "scene", sceneName: key.target} : null;
    if (key.action === "visibility") {
      const i = item(key, s);
      return i && typeof i.enabled === "boolean" ? {action: "visibility", sceneName: i.sceneName, sceneItemId: i.sceneItemId, targetId: i.targetId, enabled: !i.enabled} : null;
    }
    const i = input(key, s);
    if (key.action === "mute") return i && typeof i.muted === "boolean" ? {action: "mute", inputName: i.inputName, muted: !i.muted} : null;
    if (!i || !i.mediaState || i.mediaState === "OBS_MEDIA_STATE_NONE" || i.mediaState === "OBS_MEDIA_STATE_ERROR") return null;
    if (key.action === "media-pause" && i.mediaState !== "OBS_MEDIA_STATE_PLAYING") return null;
    if (key.action === "media-play" && i.mediaState === "OBS_MEDIA_STATE_PLAYING") return null;
    return {action: key.action, inputName: i.inputName};
  }
  function active(key, s) {
    if (!s?.connected || !key) return false;
    if (key.action === "scene") return s.currentScene === key.target;
    if (key.action === "visibility") return item(key, s)?.enabled === true;
    if (key.action === "mute") return input(key, s)?.muted === true;
    const media = input(key, s)?.mediaState;
    return key.action === "media-pause" ? media === "OBS_MEDIA_STATE_PAUSED" : key.action === "media-stop" ? media === "OBS_MEDIA_STATE_STOPPED" : media === "OBS_MEDIA_STATE_PLAYING";
  }
  function defaults(s) {
    const keys = [];
    const make = (action, target, color, icon, extra = {}) => ({action, target, color, icon, label: "", ...extra});
    (s.scenes || []).slice(0, 4).forEach(name => keys.push(make("scene", name, "amber", "scene")));
    // Put microphones before media and visual sources; only actual OBS capabilities create keys.
    (s.inputs || []).filter(i => typeof i.muted === "boolean").slice(0, 3).forEach(i => keys.push(make("mute", i.inputName, "green", "mic")));
    (s.items || []).slice().sort((a, b) => Number(b.sceneName === s.currentScene) - Number(a.sceneName === s.currentScene)).slice(0, 4).forEach(i => keys.push(make("visibility", "", "blue", "eye", {sceneName: i.sceneName, sceneItemId: i.sceneItemId, targetId: i.targetId})));
    const media = (s.inputs || []).find(i => i.mediaState && !["OBS_MEDIA_STATE_NONE", "OBS_MEDIA_STATE_ERROR"].includes(i.mediaState));
    if (media) keys.push(make("media-restart", media.inputName, "purple", "restart"));
    return Array.from({length: 12}, (_, i) => keys[i] || null);
  }
  window.OBSDeckModel = {valid, payload, active, defaults};
  const root = document.querySelector("#deck-dashboard");
  if (!root) return;
  const $ = id => document.querySelector("#" + id);
  const storageKey = "obs-controller.deck.v1";
  let state = null, layout = null, editing = false, busy = false, polling = false, selected = 0, epoch = 0;
  try { const saved = JSON.parse(localStorage.getItem(storageKey)); if (Array.isArray(saved) && saved.length === 12 && saved.every(k => k === null || valid(k))) layout = saved; } catch (_) { /* Invalid/unavailable storage: rebuild from OBS. */ }
  const feedback = value => { $("deck-feedback").textContent = value; };
  const connection = online => {
    const el = document.querySelector(".connection-status");
    el.textContent = online ? "● OBS bağlı" : "● OBS bağlı değil";
    el.classList.toggle("is-connected", online);
  };
  function persist() {
    try { localStorage.setItem(storageKey, JSON.stringify(layout)); return true; }
    catch (_) { feedback("Düzen bu oturumda hazır; tarayıcı depolaması kullanılamıyor."); return false; }
  }
  function name(key) { return key.label || (key.action === "visibility" ? item(key, state)?.sourceName || "Kaynak bulunamadı" : key.target); }
  function status(key) {
    if (!state?.connected) return "BAĞLANTI YOK";
    if (key.action === "scene") return active(key, state) ? "PROGRAMDA" : "SAHNEYE GEÇ";
    if (key.action === "visibility") return !item(key, state) ? "HEDEF YOK" : active(key, state) ? "GÖRÜNÜR" : "GİZLİ";
    if (key.action === "mute") return typeof input(key, state)?.muted !== "boolean" ? "HEDEF YOK" : active(key, state) ? "SUSTURULDU" : "SES AÇIK";
    return ({OBS_MEDIA_STATE_PLAYING: "OYNATILIYOR", OBS_MEDIA_STATE_PAUSED: "DURAKLATILDI", OBS_MEDIA_STATE_STOPPED: "DURDU", OBS_MEDIA_STATE_ENDED: "SONLANDI"})[input(key, state)?.mediaState] || "HAZIR DEĞİL";
  }
  function render() {
    const online = state?.connected === true;
    connection(online);
    $("deck-program").textContent = online ? state.currentScene || "—" : "Bağlantı bekleniyor";
    $("deck-online").textContent = online ? "● OBS ile bağlı" : "● Kontroller kilitli";
    $("deck-online").classList.toggle("online", online);
    const keys = layout || Array(12).fill(null);
    const ready = keys.filter(k => payload(k, state)).length;
    $("deck-ready").textContent = `${ready} / 12`;
    $("deck-count").textContent = `${keys.filter(Boolean).length} / 12 tuş atanmış`;
    $("deck-grid").setAttribute("aria-busy", String(busy || !layout));
    $("deck-edit").textContent = editing ? "✓ Düzenlemeyi bitir" : "Düzeni düzenle";
    $("deck-edit").setAttribute("aria-pressed", String(editing));
    $("deck-edit").disabled = busy || !layout;
    $("deck-reset").disabled = !online || busy || !editing;
    $("deck-mode-note").textContent = editing ? "Düzenleme modu · Bir yuvaya dokunarak düzenleyin" : "Gerçek OBS durumu · Komutlar tekrar gönderilmez";
    root.classList.toggle("is-editing", editing);
    // Keep each DOM button stable so polling never drops keyboard focus.
    keys.forEach((key, index) => {
      const button = $("deck-grid").children[index];
      button.className = `deck-key ${key ? "deck-color-" + key.color : "deck-empty"}`;
      button.classList.toggle("is-active", active(key, state));
      button.disabled = busy || (!editing && !payload(key, state));
      button.setAttribute("aria-pressed", String(active(key, state)));
      const label = key ? name(key) : "Tuş ekle";
      button.setAttribute("aria-label", `${index + 1}. ${label}${key ? ": " + status(key) : ""}${editing ? ", düzenle" : ""}`);
      button.title = editing ? "Bu tuşu düzenle" : !online ? "OBS bağlantısı bekleniyor" : key && !payload(key, state) ? "OBS hedefi eksik veya bu işlem mevcut durumda kullanılamıyor" : label;
      button.querySelector(".deck-key-number").textContent = String(index + 1).padStart(2, "0");
      button.querySelector(".deck-key-type").textContent = key ? titles[key.action] : "BOŞ YUVA";
      button.querySelector(".deck-key-icon").textContent = key ? icons[key.icon] : "+";
      button.querySelector(".deck-key-label").textContent = label;
      button.querySelector(".deck-key-status").textContent = key ? status(key) : editing ? "DOKUN VE ATA" : "DÜZENLEME MODUNDA EKLE";
    });
  }
  function log(label, ok) {
    const list = $("deck-log");
    if (!list.dataset.started) { list.replaceChildren(); list.dataset.started = "true"; }
    const li = document.createElement("li"), when = document.createElement("time"), message = document.createElement("span");
    when.textContent = new Date().toLocaleTimeString("tr-TR", {hour: "2-digit", minute: "2-digit"});
    message.textContent = label; li.className = ok ? "success" : "failure"; li.append(when, message); list.prepend(li);
    while (list.children.length > 3) list.lastElementChild.remove();
  }
  async function execute(index) {
    if (editing) { openEditor(index); return; }
    const key = layout?.[index], command = payload(key, state);
    if (busy || !command) return;
    busy = true; epoch++; render(); feedback("İşlem OBS’den doğrulanıyor…");
    const label = name(key);
    try {
      const response = await fetch("/api/deck/action", {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify(command), signal: AbortSignal.timeout(12000)});
      if (!response.ok) throw new Error(response.status === 409 ? "OBS hedefi veya durumu değişti; listeyi yenileyin." : response.status === 403 ? "İstek reddedildi; sayfayı yenileyin." : "İşlem doğrulanamadı. Komut tekrar gönderilmedi.");
      const next = await response.json();
      if (next?.connected !== true) throw new Error("OBS bağlantısı yok; işlem doğrulanamadı.");
      state = next; feedback(`${label} · OBS durumu doğrulandı.`); log(label, true);
    } catch (err) {
      state = null; feedback(err.name === "TimeoutError" ? "Yanıt zaman aşımı; işlem sonucu bilinmiyor. Komut tekrarlanmadı." : err.message || "Bağlantı hatası; işlem sonucu bilinmiyor."); log(label + " · doğrulanamadı", false);
    } finally { busy = false; render(); }
  }
  for (let index = 0; index < 12; index++) {
    const button = document.createElement("button"); button.type = "button"; button.dataset.deckSlot = String(index);
    for (const cls of ["deck-key-number", "deck-key-type", "deck-key-icon", "deck-key-label", "deck-key-status"]) { const el = document.createElement("span"); el.className = cls; if (cls === "deck-key-icon") el.setAttribute("aria-hidden", "true"); button.append(el); }
    button.addEventListener("click", () => execute(index)); $("deck-grid").append(button);
  }
  function targets() {
    const action = $("deck-action").value;
    if (action === "scene") return (state?.scenes || []).map(target => ({label: target, value: {target}}));
    if (action === "visibility") return (state?.items || []).map(i => ({label: `${i.sceneName} / ${i.sourceName}`, value: {sceneName: i.sceneName, sceneItemId: i.sceneItemId, targetId: i.targetId}}));
    return (state?.inputs || []).filter(i => action === "mute" ? typeof i.muted === "boolean" : i.mediaState && !["OBS_MEDIA_STATE_NONE", "OBS_MEDIA_STATE_ERROR"].includes(i.mediaState)).map(i => ({label: i.inputName, value: {target: i.inputName}}));
  }
  function fillTargets(key) {
    const select = $("deck-target"); select.replaceChildren();
    const list = targets();
    list.forEach(t => { const option = document.createElement("option"); option.value = JSON.stringify(t.value); option.textContent = t.label; select.append(option); });
    if (key) {
      const t = list.find(t => key.action === "visibility" ? typeof key.targetId === "string" && key.targetId.length > 0 && t.value.targetId === key.targetId && t.value.sceneName === key.sceneName && t.value.sceneItemId === key.sceneItemId : t.value.target === key.target);
      // Native selects otherwise silently choose the first remaining OBS target.
      select.selectedIndex = t ? list.indexOf(t) : -1;
    }
    select.disabled = list.length === 0;
    updateEditorTarget(key && select.selectedIndex === -1 ? "Kayıtlı OBS hedefi bulunamadı veya kimliği değişti. Kaydetmek için açıkça yeni bir hedef seçin." : "");
  }
  function updateEditorTarget(warning = "") {
    const select = $("deck-target");
    const available = state?.connected && targets().some(t => JSON.stringify(t.value) === select.value);
    $("deck-save").disabled = !available;
    $("deck-editor-feedback").textContent = available ? "" : warning || "Bu eylem için geçerli bir OBS hedefi seçin.";
  }
  $("deck-target").addEventListener("change", () => updateEditorTarget());
  function openEditor(index) {
    selected = index; const key = layout[index];
    $("deck-editor-title").textContent = `Tuş ${String(index + 1).padStart(2, "0")} · ${key ? "Düzenle" : "Yeni tuş"}`;
    $("deck-label").value = key?.label || ""; $("deck-action").value = key?.action || "scene"; $("deck-color").value = key?.color || "amber"; $("deck-icon").value = key?.icon || "scene";
    $("deck-position").value = String(index); $("deck-delete").disabled = !key; fillTargets(key); $("deck-editor").showModal();
  }
  for (let i = 0; i < 12; i++) { const option = document.createElement("option"); option.value = String(i); option.textContent = `Tuş ${String(i + 1).padStart(2, "0")}`; $("deck-position").append(option); }
  $("deck-action").addEventListener("change", () => {
    const a = $("deck-action").value; $("deck-color").value = a === "scene" ? "amber" : a === "visibility" ? "blue" : a === "mute" ? "green" : "purple";
    $("deck-icon").value = a === "scene" ? "scene" : a === "visibility" ? "eye" : a === "mute" ? "mic" : a.slice(6); fillTargets();
  });
  $("deck-form").addEventListener("submit", event => {
    event.preventDefault(); updateEditorTarget(); if ($("deck-save").disabled) return;
    const key = {action: $("deck-action").value, ...JSON.parse($("deck-target").value), label: $("deck-label").value.trim(), color: $("deck-color").value, icon: $("deck-icon").value};
    if (!valid(key)) return;
    const position = Number($("deck-position").value); layout[selected] = layout[position]; layout[position] = key;
    if (persist()) feedback("Tuş kaydedildi. Düzen bu tarayıcıda korunur."); $("deck-editor").close(); render();
  });
  $("deck-delete").addEventListener("click", () => { layout[selected] = null; persist(); $("deck-editor").close(); render(); });
  $("deck-close").addEventListener("click", () => $("deck-editor").close());
  $("deck-edit").addEventListener("click", () => { editing = !editing; render(); });
  $("deck-reset").addEventListener("click", () => { if (!editing || !state?.connected || !window.confirm("Özel tuş düzeni silinip gerçek OBS hedefleriyle yeniden doldurulsun mu?")) return; layout = defaults(state); persist(); render(); feedback("Deck gerçek OBS hedefleriyle dolduruldu."); });
  async function poll() {
    if (busy || polling) { setTimeout(poll, 1500); return; }
    polling = true; const version = epoch;
    try {
      const response = await fetch("/api/deck", {cache: "no-store", signal: AbortSignal.timeout(10000)});
      if (!response.ok) throw new Error("status");
      const next = await response.json();
      if (typeof next?.connected !== "boolean") throw new Error("invalid");
      if (version !== epoch) return;
      state = next;
      if (!layout && state.connected) { layout = defaults(state); persist(); feedback("Deck hazır. Bir tuşa dokunun veya düzeni kişiselleştirin."); }
    } catch (_) { if (version === epoch) { state = null; feedback("OBS bağlantısı yok veya durum alınamadı. Kontroller kilitli; komutlar tekrar gönderilmez."); } }
    finally { polling = false; if (version === epoch) render(); setTimeout(poll, 1500); }
  }
  render(); poll();
})();
