# Ayarlar ek güvenlik doğrulaması

- Tarih/saat: 2026-11-02 00:57 +03. Linux, geçici Node ve Chromium araçları.
- Amaç: Yerel/LAN arayüz ayrımı, Origin ve bozuk yapılandırma regresyonları dahil tam test grubu.
- Komut: `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./...` — çıkış 0, tüm paketler başarılı; kök paket önbellekten.
- Gerçek OBS kullanılmadı; testler simüle edildi. Gerçek OBS kabulü, fiziksel cihaz ve ayar penceresinin gerçek tarayıcı etkileşimi doğrulanmadı. Kayıt testleri etkinleştirilmedi.
