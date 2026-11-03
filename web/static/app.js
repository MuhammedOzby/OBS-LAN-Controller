"use strict";

(() => {
  // Deck owns its serialized inventory polling and connection indicator.
  if (document.querySelector("#deck-dashboard")) return;
  const dashboard = document.querySelector("#control-dashboard");
  const slots = [...document.querySelectorAll("[data-scene-slot]")];
  const connection = document.querySelector(".connection-status");
  let state = null;
  let sceneBusy = false;
  let statusBusy = false;
  const text = (id, value) => {
    const target = document.querySelector("#" + id);
    if (target) target.textContent = value;
  };
  const finite = value => typeof value === "number" && Number.isFinite(value) && value >= 0;
  const metric = (value, unit, decimals = 1) => finite(value) ? value.toFixed(decimals) + unit : "—";
  const duration = value => {
    if (!finite(value)) return "—";
    const seconds = Math.floor(value / 1000);
    return [Math.floor(seconds / 3600), Math.floor(seconds / 60) % 60, seconds % 60].map(n => String(n).padStart(2, "0")).join(":");
  };
  function renderSlots() {
    slots.forEach((button, index) => {
      const name = state?.connected ? (state.favorites || state.scenes)?.[index] : null;
      const valid = typeof name === "string" && name.length > 0 && state.scenes?.includes(name);
      button.textContent = valid ? name : name ? `${name} (kullanılamıyor)` : `Sahne ${index + 1} yok`;
      button.disabled = !valid || sceneBusy || statusBusy;
      button.setAttribute("aria-pressed", String(valid && name === state.currentScene));
      button.title = valid ? name : "Bu sahne yuvası boş";
    });
  }
  function render(next) {
    state = next;
    const online = next?.connected === true;
    connection.textContent = online ? "● OBS bağlı" : next?.reconnecting ? "● OBS yeniden bağlanıyor" : "● OBS bağlı değil";
    connection.title = !online && typeof next?.connectionError === "string" ? next.connectionError : "";
    connection.classList.toggle("is-connected", online);
    window.OBSBroadcast?.render(next);
    window.OBSAudio?.connection(online);
    if (!dashboard) return;
    const live = online && next.streamActive === true;
    text("stream-state", live ? "Canlı" : online ? "Yayın yapılmıyor" : "Yayın durumu bilinmiyor");
    document.querySelector("#stream-state").classList.toggle("is-live", live);
    text("stream-duration", live ? duration(next.duration) : "—");
    text("current-scene", online && typeof next.currentScene === "string" ? next.currentScene : "Sahne bekleniyor");
    text("obs-fps", online ? metric(next.fps, "", 2) : "—");
    text("obs-memory", online ? metric(next.memory, " MB", 0) : "—");
    text("obs-render", online ? metric(next.renderTime, " ms") : "—");
    text("obs-skips", online && finite(next.renderSkipped) && finite(next.encodeSkipped) ? `${next.renderSkipped} / ${next.encodeSkipped}` : "—");
    text("stat-cpu", online ? metric(next.cpu, " %") : "—");
    text("stat-network", live ? metric(next.bitrateKbps, " kb/s", 0) : "—");
    text("stat-drop", live && finite(next.dropped) && finite(next.total) ? `${next.dropped} / ${next.total}` : "—");
    // Report only measured congestion; no invented overall quality grade.
    text("stat-quality", live && finite(next.congestion) ? `Tıkanıklık ${Math.round(next.congestion * 100)} %` : "Bilinmiyor");

    if (!sceneBusy) text("scene-feedback", online ? "Program sahnesini değiştirmek için dokunun" : next?.connectionError || "OBS bağlantısı bekleniyor");
    renderSlots();
  }
  async function poll() {
    statusBusy = true;
    try {
      const response = await fetch("/api/status?preview=0", {cache: "no-store"});
      if (!response.ok) throw new Error("status");
      const next = await response.json();
      if (!next || typeof next.connected !== "boolean") throw new Error("invalid status");
      render(next);
    } catch (_) {
      render(null);
    } finally {
      statusBusy = false;
      renderSlots();
      setTimeout(poll, 1000);
    }
  }
  slots.forEach((button, index) => button.addEventListener("click", async () => {
    const sceneName = (state?.favorites || state?.scenes)?.[index];
    if (sceneBusy || statusBusy || !state?.connected || button.disabled || typeof sceneName !== "string") return;
    sceneBusy = true;
    renderSlots();
    text("scene-feedback", "Sahne değişimi doğrulanıyor…");
    try {
      const response = await fetch("/api/scene", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({sceneName})
      });
      if (response.status === 204) {
        // Do not optimistically mark active: the next status poll supplies program state.
        text("scene-feedback", "Sahne değişimi doğrulandı");
      } else if (response.status === 403) {
        text("scene-feedback", "İstek reddedildi; sayfayı yenileyin");
      } else if (response.status === 400) {
        text("scene-feedback", "Sahne bulunamadı; liste yenilenecek");
      } else {
        render(null);
        text("scene-feedback", "OBS bağlantısı yok veya sahne doğrulanamadı");
      }
    } catch (_) {
      render(null);
      text("scene-feedback", "Bağlantı hatası; işlem doğrulanamadı");
    } finally {
      sceneBusy = false;
      renderSlots();
    }
  }));
  const editor = document.querySelector("#favorites-editor");
  const edit = document.querySelector("#favorites-edit");
  if (dashboard && editor && edit) {
    const form = document.querySelector("#favorites-form");
    const selects = [0, 1, 2].map(i => document.querySelector(`#favorite-${i}`));
    let collection = "";
    edit.addEventListener("click", async () => {
      text("favorites-feedback", "Sahneler okunuyor…");
      editor.showModal();
      try {
        const response = await fetch("/api/favorites", {cache: "no-store"});
        if (!response.ok) throw new Error();
        const data = await response.json(); collection = data.collection;
        selects.forEach((select, i) => {
          select.replaceChildren();
          const add = (value, label) => { const option = document.createElement("option"); option.value = value; option.textContent = label; select.append(option); };
          add("", "Boş");
          data.scenes.forEach(name => add(name, name));
          if (data.slots[i] && !data.scenes.includes(data.slots[i])) add(data.slots[i], `${data.slots[i]} (kullanılamıyor)`);
          select.value = data.slots[i] || "";
        });
        text("favorites-feedback", "");
      } catch (_) { collection = ""; text("favorites-feedback", "Sahneler okunamadı."); }
    });
    document.querySelector("#favorites-close").addEventListener("click", () => editor.close());
    form.addEventListener("submit", async event => {
      event.preventDefault();
      if (!collection) return;
      try {
        const response = await fetch("/api/favorites", {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify({collection, slots: selects.map(s => s.value)})});
        if (!response.ok) throw new Error();
        editor.close();
        if (state) state.favorites = selects.map(s => s.value);
        renderSlots();
      } catch (_) { text("favorites-feedback", "Kaydedilemedi; seçimi ve OBS bağlantısını kontrol edin."); }
    });
  }
  poll();
})();
