# QR son doğrulama

- Tarih: 2 Kasım 2026 14:35:49–14:36:05 +03, Linux, geçici Node ve headless Chrome.
- Amaç: Son kullanım saati ve geçersiz kodda elle giriş düzeltmelerinden sonra regresyon.
- `PATH=/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test -race ./...`: çıkış 0, tüm paketler geçti.
- `go vet ./...`: çıkış 0. `go build -o /tmp/opencode/obs-lan-controller-qr .`: çıkış 0.
- OBS simüle edildi veya kullanılmadı; gerçek OBS, fiziksel telefon/kamera ve PWA doğrulanmadı. Kullanıcı sunucusuna dokunulmadı.
