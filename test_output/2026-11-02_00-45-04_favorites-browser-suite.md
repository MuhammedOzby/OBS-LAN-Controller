# Favori sahneler tarayıcı ve regresyon testleri

- Tarih: 2026-11-02 00:45:04 +03 (+0300).
- Amaç: üç favori yuvası, düzenleyici ve mevcut regresyonları doğrulamak.
- Komut: `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./...`.
- Sonuç: başarılı, çıkış kodu 0; Node davranış testi ve headless Chromium düzenleyici testi dahil; Lib sonuçları önbellekten.
- OBS ve HTTP verileri simüle edildi; gerçek OBS, fiziksel cihaz ve gerçek yayın/kayıt değiştirilmedi.
