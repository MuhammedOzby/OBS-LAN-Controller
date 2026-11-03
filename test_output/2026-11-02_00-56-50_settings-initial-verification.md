# Ayarlar ilk doğrulama

- Tarih/saat: 2026-11-02 00:56:50 +03. Linux, geçici Node ve Chromium araçları.
- Amaç: İlk kurulum/ayar değişiklikleri ve mevcut regresyonları doğrulamak.
- Komut: `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./...` — çıkış 0, tüm paketler başarılı.
- Komutlar: `go vet ./...` ve `go build -o /tmp/opencode/obs-lan-controller-settings-check .` — ikisi de çıkış 0.
- Testler simülasyon kullanır; gerçek OBS'ye bağlanılmadı, kayıt testleri çalıştırılmadı. Fiziksel cihaz, gerçek OBS ve ayar penceresinin gerçek tarayıcıda etkileşimi doğrulanmadı.
