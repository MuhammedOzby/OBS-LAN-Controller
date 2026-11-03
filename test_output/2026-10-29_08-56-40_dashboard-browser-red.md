# Gerçek OBS tarayıcı — ilk bütünleşik kontrol

Tarih: 2026-10-29, UTC. Linux Chromium / http://127.0.0.1:8080.

Adımlar: gerçek OBS configuration.json ile sunucu; sayfa yenileme ve DOM metrikleri; Emulation.setDeviceMetricsOverride ile 844×390, 932×430, 390×844, 1280×720.

İlk tarayıcı denemesi CDP bağlantı kopmasıyla engellendi (`no close frame received or sent`); ensure_real_tab ile toparlandı.

Gerçek OBS bağlı, Canlı, JPEG görüntü naturalWidth=960, FPS 30.00, CPU 0.8%, AĞ 7463 kb/s, DROP 0/26832 okundu. Yalnız bir Scene mevcut; diğer iki yuva disabled. Sahne veya yayın değiştirilmedi.

Sonuç: masaüstü ve iki yatay boyutta taşma yok, önizleme oranı ~1.7778. 390×844 dikey boyutta stats-table içeriği minimum genişliği nedeniyle scene-column/control-dashboard 6px taşıyor (358px clientWidth, 364px scrollWidth). Düzeltme gerekli. Gerçek görüntü yüklemesi doğrulandı; fiziksel telefon/uzun süreli performans testi değildir.
