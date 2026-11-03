# Gerçek OBS bağlantı testi

Tarih: 2026-10-29T08:38:44.627065+00:00

Hedef: 192.168.88.254:4455

Ortam: Linux, Python websockets.

Adımlar: configuration.json üzerinden bağlantı; OBS WebSocket v5 Hello/Identify ve challenge-response kimlik doğrulaması; salt okunur GetVersion ve GetStats istekleri.

- WebSocket bağlantısı kuruldu.
- OBS kimlik doğrulaması başarılı.
- OBS Studio: 33.0.0; obs-websocket: 5.7.6
- GetStats başarılı; aktif FPS: 30.0; CPU: %0.3

Çıkış kodu: 0

Gerçek OBS kullanıldı; uygulamanın Go OBS istemcisi henüz uygulanmadı. Sahne/ses/yayın/kayıt durumunda değişiklik yapılmadı. Kalıcı bağlantı ve yeniden bağlanma test edilmedi.
