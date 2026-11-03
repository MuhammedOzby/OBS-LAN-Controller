"use strict";

(() => {
  if (!document.querySelector("#control-dashboard")) return;
  const video = document.querySelector("#program-video");
  const placeholder = document.querySelector("#preview-placeholder");
  const status = document.querySelector("#preview-status");
  const fps = document.querySelector("#preview-fps");
  const retry = document.querySelector("#preview-retry");
  if (!video || !placeholder || !status || !fps || !retry) return;

  const MAX_ATTEMPTS = 5;
  let current = null;
  let attempts = 0;
  let retryTimer = null;
  let pageActive = true;
  const visible = () => pageActive && !document.hidden;
  const valid = ctx => current === ctx && visible();
  const message = value => { status.textContent = value; };
  const unavailable = value => {
    video.hidden = true;
    placeholder.hidden = false;
    placeholder.textContent = value;
    fps.textContent = "—";
  };
  video.muted = true;
  video.autoplay = true;
  video.playsInline = true;

  function closeSession(id) {
    if (!id) return;
    // Keepalive also covers navigation and a session answer arriving after hide.
    void fetch("/api/preview/session/" + encodeURIComponent(id), {
      method: "DELETE", credentials: "same-origin", cache: "no-store", keepalive: true
    }).catch(() => {});
  }
  function later(ctx, fn, ms) {
    const timer = setTimeout(() => {
      ctx.timers.delete(timer);
      if (valid(ctx)) fn();
    }, ms);
    ctx.timers.add(timer);
    return timer;
  }
  function clear(ctx, timer) {
    clearTimeout(timer);
    ctx.timers.delete(timer);
  }
  function stop() {
    clearTimeout(retryTimer);
    retryTimer = null;
    const ctx = current;
    current = null;
    if (ctx) {
      ctx.cancel.abort();
      ctx.requests.forEach(controller => controller.abort());
      ctx.timers.forEach(clearTimeout);
      ctx.pc.ontrack = null;
      ctx.pc.onconnectionstatechange = null;
      ctx.pc.close();
      closeSession(ctx.id);
    }
    if (video.srcObject) video.srcObject.getTracks().forEach(track => track.stop?.());
    video.pause();
    video.srcObject = null;
    unavailable("Önizleme durduruldu");
  }
  function fail(ctx, error) {
    if (!valid(ctx)) return;
    const denied = error?.denied === true;
    stop();
    const exhausted = denied || attempts >= MAX_ATTEMPTS;
    const reason = error?.message || "Önizleme bağlantısı kesildi";
    message(exhausted ? reason + " · " + (denied ? "Sayfayı yenileyin" : "Yeniden deneyebilirsiniz") : reason + " · Yeniden bağlanılıyor");
    unavailable(reason);
    retry.hidden = denied || !exhausted;
    if (!exhausted && visible()) {
      retryTimer = setTimeout(() => {
        retryTimer = null;
        if (visible()) void start();
      }, Math.min(1000 * 2 ** (attempts - 1), 8000));
    }
  }
  async function request(ctx, url, options = {}) {
    const controller = new AbortController();
    ctx.requests.add(controller);
    const timer = later(ctx, () => controller.abort(), 5000);
    try {
      const response = await fetch(url, {...options, credentials: "same-origin", cache: "no-store", signal: controller.signal});
      if (!response.ok) {
        const error = new Error(response.status === 401 || response.status === 403 ? "Önizleme erişimi reddedildi" : "Önizleme servisine ulaşılamıyor");
        error.denied = response.status === 401 || response.status === 403;
        throw error;
      }
      return await response.json();
    } finally {
      clear(ctx, timer);
      ctx.requests.delete(controller);
    }
  }
  function gather(ctx) {
    if (ctx.pc.iceGatheringState === "complete") return Promise.resolve();
    return new Promise((resolve, reject) => {
      const finish = error => {
        clear(ctx, timer);
        ctx.pc.removeEventListener("icegatheringstatechange", changed);
        ctx.cancel.signal.removeEventListener("abort", aborted);
        if (error) reject(error); else resolve();
      };
      const changed = () => { if (ctx.pc.iceGatheringState === "complete") finish(); };
      const aborted = () => finish(new Error("Önizleme durduruldu"));
      const timer = later(ctx, () => finish(new Error("Önizleme ICE zaman aşımı")), 5000);
      ctx.pc.addEventListener("icegatheringstatechange", changed);
      ctx.cancel.signal.addEventListener("abort", aborted, {once: true});
      changed();
    });
  }
  async function monitor(ctx) {
    try {
      const source = await request(ctx, "/api/preview/status");
      if (!valid(ctx)) return;
      if (!source || typeof source.connected !== "boolean") throw new Error("Geçersiz önizleme durumu");
      // ICE and the first encoder keyframe arrive after the SDP response.
      // A pending source is not a failure during the bounded startup window.
      if ((!source.connected || !source.streaming) && performance.now() - ctx.started > 12000) throw new Error(source.error || "Önizleme kaynağı bekleniyor");
      // Browser receive statistics, never OBS FPS or the encoder's target FPS.
      const stats = await ctx.pc.getStats();
      if (!valid(ctx)) return;
      let measured = null;
      stats.forEach(report => {
        if (report.type !== "inbound-rtp" || (report.kind || report.mediaType) !== "video" || report.isRemote) return;
        if (typeof report.framesDecoded === "number") {
          if (!ctx.sample || report.framesDecoded > ctx.sample.frames) ctx.lastFrame = performance.now();
          const dt = ctx.sample ? report.timestamp - ctx.sample.time : 0;
          if (ctx.sample && dt > 0 && report.framesDecoded >= ctx.sample.frames) measured = (report.framesDecoded - ctx.sample.frames) * 1000 / dt;
          ctx.sample = {frames: report.framesDecoded, time: report.timestamp};
        }
        if (typeof report.framesPerSecond === "number") measured = report.framesPerSecond;
      });
      fps.textContent = !video.hidden && Number.isFinite(measured) && measured >= 0 ? measured.toFixed(1) : "—";
      if (performance.now() - ctx.lastFrame > 12000) throw new Error("Önizleme video akışı durdu");
      if (!video.hidden) message("Önizleme bağlı · WebRTC");
    } catch (error) {
      fail(ctx, error);
      return;
    }
    if (valid(ctx)) later(ctx, () => { void monitor(ctx); }, 1000);
  }
  async function start() {
    if (!visible() || current) return;
    retry.hidden = true;
    if (typeof RTCPeerConnection !== "function") {
      message("Bu tarayıcı WebRTC önizlemesini desteklemiyor");
      unavailable("WebRTC desteklenmiyor; güvenli bağlantı veya güncel tarayıcı kullanın");
      return;
    }
    attempts++;
    let ctx;
    try {
      ctx = {pc: new RTCPeerConnection({iceServers: []}), id: null, requests: new Set(), timers: new Set(), cancel: new AbortController(), started: performance.now(), lastFrame: performance.now(), sample: null};
      current = ctx;
      message("Önizleme bağlanıyor · WebRTC");
      unavailable("Program video akışı bekleniyor");
      ctx.pc.addTransceiver("video", {direction: "recvonly"});
      ctx.pc.ontrack = event => {
        if (!valid(ctx) || event.track.kind !== "video") return;
        video.srcObject = event.streams[0] || new MediaStream([event.track]);
        event.track.onended = () => fail(ctx, new Error("Önizleme video izi kapandı"));
        void video.play().catch(() => fail(ctx, new Error("Önizleme otomatik oynatılamadı")));
      };
      ctx.pc.onconnectionstatechange = () => {
        if (["failed", "disconnected", "closed"].includes(ctx.pc.connectionState)) fail(ctx, new Error("Önizleme bağlantısı kesildi"));
      };
      const offer = await ctx.pc.createOffer();
      if (!valid(ctx)) return;
      await ctx.pc.setLocalDescription(offer);
      if (!valid(ctx)) return;
      await gather(ctx);
      if (!valid(ctx)) return;
      const answer = await request(ctx, "/api/preview/offer", {
        method: "POST", headers: {"Content-Type": "application/json"},
        body: JSON.stringify({type: ctx.pc.localDescription.type, sdp: ctx.pc.localDescription.sdp})
      });
      // Capture the server ID before validating SDP, including stale responses.
      const id = typeof answer?.id === "string" && answer.id.length > 0 ? answer.id : null;
      if (!valid(ctx)) { closeSession(id); return; }
      ctx.id = id;
      if (!id || answer.type !== "answer" || typeof answer.sdp !== "string" || !answer.sdp) throw new Error("Geçersiz önizleme yanıtı");
      await ctx.pc.setRemoteDescription({type: answer.type, sdp: answer.sdp});
      if (valid(ctx)) {
        later(ctx, () => { if (video.hidden) fail(ctx, new Error("Önizleme video zaman aşımı")); }, 12000);
        void monitor(ctx);
      }
    } catch (error) {
      if (ctx) fail(ctx, error);
      else {
        message("WebRTC başlatılamadı; tarayıcı bağlantısını kontrol edin");
        unavailable("Önizleme kullanılamıyor");
        retry.hidden = false;
      }
    }
  }
  video.addEventListener("playing", () => {
    if (!current || !visible() || !video.srcObject) return;
    video.hidden = false;
    placeholder.hidden = true;
    current.lastFrame = performance.now();
    message("Önizleme bağlı · WebRTC");
  });
  retry.addEventListener("click", () => {
    stop();
    attempts = 0;
    if (visible()) void start();
  });
  document.addEventListener("visibilitychange", () => {
    stop();
    if (visible()) { attempts = 0; void start(); }
    else message("Önizleme duraklatıldı · Sayfa görünür değil");
  });
  window.addEventListener("pagehide", () => { pageActive = false; stop(); });
  window.addEventListener("pageshow", () => {
    pageActive = true;
    if (visible() && !current) { attempts = 0; void start(); }
  });
  if (visible()) void start();
})();
