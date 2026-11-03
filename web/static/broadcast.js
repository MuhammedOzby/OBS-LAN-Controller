"use strict";

(() => {
  const root = document.querySelector("#broadcast-dashboard");
  if (!root) return;
  const buttons = [...document.querySelectorAll("[data-broadcast-action]")];

  const labels = {
    "start-stream": "Yayını başlat",
    "stop-stream": "Yayını durdur",
    "start-record": "Kaydı başlat",
    "stop-record": "Kaydı durdur",
    "pause-record": "Kaydı duraklat",
    "resume-record": "Kayda devam et"
  };
  let state = null;
  let busy = false;

  let logCount = 0;
  const text = (id, value) => { document.querySelector("#" + id).textContent = value; };
  const number = (value, suffix = "", digits = 1) => typeof value === "number" && Number.isFinite(value) && value >= 0 ? value.toFixed(digits) + suffix : "—";
  const bytes = value => typeof value === "number" && Number.isFinite(value) && value >= 0 ? (value / 1048576).toFixed(1) + " MB" : "—";
  const clock = value => {
    if (typeof value !== "number" || !Number.isFinite(value) || value < 0) return "—";
    const s = Math.floor(value / 1000);
    return [Math.floor(s / 3600), Math.floor(s / 60) % 60, s % 60].map(n => String(n).padStart(2, "0")).join(":");
  };
  function allowed(action) {
    if (!state?.connected || busy) return false;
    if (action === "start-stream") return !state.streamActive && !state.reconnecting;
    if (action === "stop-stream") return state.streamActive === true;
    if (!state.recordAvailable) return false;
    if (action === "start-record") return !state.recordActive;
    if (action === "stop-record") return state.recordActive === true;
    if (state.recordPauseSupported === false) return false;
    if (action === "pause-record") return state.recordActive === true && !state.recordPaused;
    return action === "resume-record" && state.recordActive === true && state.recordPaused === true;
  }
  function controls() {
    buttons.forEach(button => {
      const action = button.dataset.broadcastAction;
      button.disabled = !allowed(action);
      const unsupported = (action === "pause-record" || action === "resume-record") && state?.recordPauseSupported === false;
      const hint = unsupported ? "Bu OBS ayarında duraklatma/devam desteklenmiyor: kayıt, yayın encoder’ını paylaşıyor (Yayınla aynı)." : !state?.connected ? "OBS bağlantısı bekleniyor." : busy ? "İşlem tamamlanana kadar bekleyin." : button.disabled ? "Mevcut OBS durumunda bu işlem kullanılamaz." : labels[action];
      button.title = hint;
      if (button.parentElement?.classList.contains("record-action-hint")) {
        button.parentElement.title = hint;
        button.parentElement.setAttribute("aria-label", labels[action] + ": " + hint);
      }
    });
    root.setAttribute("aria-busy", String(busy));
  }
  function render(next) {
    // In-flight polls must not overwrite the command's authoritative readback.
    if (busy) return;
    state = next;
    const online = next?.connected === true;
    const live = online && next.streamActive === true;
    const record = online && next.recordAvailable === true;
    const recording = record && next.recordActive === true;
    text("manage-stream-state", !online ? "BİLİNMİYOR" : next.reconnecting ? "YENİDEN BAĞLANIYOR" : live ? "CANLI" : "YAYIN KAPALI");
    text("manage-record-state", !record ? "BİLİNMİYOR" : recording ? next.recordPaused ? "DURAKLATILDI" : "KAYITTA" : "KAYIT KAPALI");
    document.querySelector("#manage-stream-state").classList.toggle("is-active", live);
    document.querySelector("#manage-record-state").classList.toggle("is-active", recording);
    text("manage-stream-duration", online ? clock(live ? next.duration : 0) : "—");
    text("manage-record-duration", record ? clock(recording ? next.recordDuration : 0) : "—");
    text("manage-bitrate", live ? number(next.bitrateKbps, " kb/s", 0) : "—");
    text("manage-stream-bytes", online ? bytes(next.streamBytes) : "—");
    text("manage-record-bytes", record ? bytes(next.recordBytes) : "—");
    text("manage-dropped", live ? number(next.dropped, " kare", 0) : "—");
    text("manage-record-mode", !record ? "Alınamıyor" : recording ? next.recordPaused ? "Beklemede" : "Yazılıyor" : "Hazır");
    text("record-pause-note", next?.recordPauseSupported === false ? "Duraklatma kapalı: OBS kayıt kalitesinde ‘Yayınla aynı’ / paylaşılan encoder kullanılıyor. Ayrı kayıt encoder’ı seçin." : "Duraklatma OBS / encoder desteğine bağlı");
    text("manage-fps", online ? number(next.fps, "", 2) : "—");
    text("manage-cpu", online ? number(next.cpu, " %") : "—");
    text("manage-render", online ? number(next.renderTime, " ms") : "—");
    text("manage-skips", online ? number(next.renderSkipped, "", 0) + " / " + number(next.encodeSkipped, "", 0) : "—");
    text("manage-scene", online && typeof next.currentScene === "string" ? next.currentScene : "—");
    text("broadcast-health-note", !online ? "OBS bağlantısı yok. Kritik kontroller kilitli." : next.reconnecting ? "OBS yayın servisine yeniden bağlanıyor." : "OBS verileri yaklaşık saniyede bir yenilenir.");

    controls();
  }
  function log(message) {
    const list = document.querySelector("#broadcast-log");
    if (logCount === 0) list.textContent = "";
    const item = document.createElement("li");
    item.textContent = new Date().toLocaleTimeString("tr-TR", {hour: "2-digit", minute: "2-digit", second: "2-digit"}) + " · " + message;
    list.prepend(item);
    logCount++;
    while (list.children.length > 4) list.children[list.children.length - 1].remove();
    text("broadcast-feedback", message);
  }
  buttons.forEach(button => button.addEventListener("click", async () => {
    const action = button.dataset.broadcastAction;
    if (!allowed(action)) return;
    busy = true;
    controls();
    text("broadcast-feedback", "OBS yanıtı ve çıkış durumu doğrulanıyor…");
    let verified = null;
    let message;
    try {
      const response = await fetch("/api/broadcast", {
        method: "POST", headers: {"Content-Type": "application/json"},
        body: JSON.stringify({action, confirmed: true}), signal: AbortSignal.timeout(12000)
      });
      if (!response.ok) {
        message = response.status === 409 ? "Durum değişti veya işlem desteklenmiyor. Yeniden kontrol edin." : response.status === 403 ? "Güvenlik kontrolü isteği reddetti. Sayfayı yenileyin." : "OBS işlemi doğrulanamadı. Tekrar göndermeden önce durumu kontrol edin.";
      } else {
        verified = await response.json();
        if (verified?.connected !== true) throw new Error("unverified");
        message = labels[action] + " · OBS durumu doğrulandı.";
      }
    } catch (_) {
      message = "Yanıt alınamadı. İşlem uygulanmış olabilir; otomatik tekrar gönderilmez.";
    } finally {
      busy = false;
      render(verified);
      log(message);
    }
  }));
  window.OBSBroadcast = {render};
})();
