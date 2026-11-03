# Başlatma kontrolü

Tarih: 2026-11-02T14:15:29.135680+03:00

Amaç: Kullanıcının elle denemesi için Controller başlatma.

Komut: `go run . serve` (arka planda). Python urllib ile loopback GET /, /api/status, /api/preview/status.

Sonuç:
- /: HTTP 200
- /api/status: HTTP 200; {"connected": true, "streamActive": false}
- /api/preview/status: HTTP 200; {"connected": true, "streaming": false}

Gerçek yapılandırma kullanıldı; simülasyon yok. Yayın/kayıt/sahne/ses komutu gönderilmedi. Video çözümü ve etkileşimli arayüz testi yapılmadı. Sunucu çalışır bırakıldı; henüz çıkış kodu yok.
