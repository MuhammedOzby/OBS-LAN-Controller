# QR: koordinatörün bağımsız doğrulaması

- Zaman: 2026-11-02, 14:40:11 +03:00 ve takip eden test süresi.
- Amaç: `/pair` QR, dördüncü sekme, güvenlik ve dar ekran düzeltmeleri sonrası bağımsız regresyon kontrolü.
- Ortam: Linux, taşınabilir Node v24.22.2 ve Chrome for Testing 156; gerçek OBS kullanılmadı.

| Komut | Çıkış | Sonuç |
| --- | --- | --- |
| `env PATH="/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH" CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test -race ./... -count=1` | 0 | Ana paket 3.174 s, Lib 14.384 s; web test dosyası yok |
| `go vet ./...` | 0 | Hata yok |
| `go build -o /tmp/opencode/obs-lan-controller-qr-reviewed .` | 0 | Geçici derleme üretildi |
| `git diff --check` | 0 | Yama/boşluk kontrolü; uygulama testi değildir |

QR kodlar bağımsız çözücü ile kontrol edildi; headless Chrome QR bağlantısından gerçek yerel/LAN soket sınırından oturum alarak kontrol sayfasına ulaştı. Diğer OBS davranışları mock verilerle test edildi. Gerçek telefon kamerası, gerçek kullanıcı LAN erişimi veya gerçek OBS komutları denenmedi. Kullanıcının çalışan Controller sunucusu ve OBS değiştirilmedi/yeniden başlatılmadı; yeni sürüme geçiş kullanıcı tarafından uygun zamanda yapılmalıdır. Oturum TTL/yeniden başlatmada iptal davranışı değiştirilmedi. Gizli kodlar/QR bağlantıları rapora yazılmadı.
