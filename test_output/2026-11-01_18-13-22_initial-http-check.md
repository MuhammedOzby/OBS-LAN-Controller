# Başlangıç HTTP kontrolü

- Tarih: 2026-11-01 18:13:22 +03:00.
- Amaç: mevcut Controller sunucusunun çalışıp çalışmadığını anlamak.
- Komut: `curl --silent --show-error --max-time 5 http://127.0.0.1:8080/api/status`
- Sonuç: bağlantı reddedildi; çıkış kodu 7. Bu adreste Controller henüz çalışmıyor.
- Bu bir gerçek yerel HTTP erişim denemesidir; OBS bağlantısı veya uygulama işlevleri doğrulanmadı. Simülasyon çalıştırılmadı.
