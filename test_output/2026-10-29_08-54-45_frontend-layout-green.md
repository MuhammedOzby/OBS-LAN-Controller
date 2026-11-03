# Frontend Chromium GREEN

UTC: 2026-10-29T08:54:45.344697+00:00

Amaç/ortam: gerçek Chromium, config okumayan localhost:18081 Go frontend harness. Gerçek OBS kullanılmadı.

Adımlar: CDP Emulation.setDeviceMetricsOverride; /, /deck, /yayin-yonetimi gezin; tüm panel/çocuk boundingClientRect ve scrollHeight/clientHeight kontrolü.

Sonuç: 15/15 sayfa-boyut senaryosu başarılı; issues=[]; document boyutu viewport ile aynı.

| Viewport | Kontrol preview en/boy | Üç sayfada taşma |
|---|---:|---|
|844×390|1.777927|yok|
|932×430|1.777889|yok|
|390×844|1.777778|yok|
|1280×720|1.777814|yok|
|1920×1080|1.777778|yok|

İkinci grup: /static/app.js gerçek kodu gerçek DOM üzerinde simüle fetch yanıtlarıyla çalıştırıldı; 844×390,932×430,390×844,1280×720. Canlı durum, süre/metrikler ve geçici PNG görüntüsü işlendi. Uzun HTML-benzeri sahne adları textContent ile düz metin kaldı; unsafeNodes=0. Tüm dört boyutta issues=[]; quality=Tıkanıklık 3 %; image=true. Bunlar sentetik fixture değerleridir, gerçek OBS ölçümü değildir.

Sınır: fiziksel cihaz/görsel screenshot incelemesi yok; gerçek OBS görüntüsü ve gerçek sahne mutasyonu doğrulanmadı. Geçici test sunucuları her embedded asset değişiminden sonra yeniden başlatıldı.
