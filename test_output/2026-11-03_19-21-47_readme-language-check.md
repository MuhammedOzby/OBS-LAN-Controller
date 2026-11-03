# README dil düzenlemesi — belge kontrolü

Tarih: 3 Kasım 2026, 19:21:47 Europe/Istanbul (UTC+03:00).

## Amaç ve kapsam

README.md İngilizce ana belge olarak düzenlendi. Önceki Türkçe gövde README-TR.md dosyasına korundu. Her iki dilde başlığın hemen altında AI desteğiyle geliştirme notu ve karşılıklı dil bağlantısı bulunur.

## Kontroller

Bağlı GitHub aracıyla güncel README, AGENTS.md ve PLAN.md okundu. Hazırlanan metinler JavaScript doğrulamasıyla kontrol edildi:

- İngilizce sürümdeki tüm fenced komut/veri blokları, yorum çevirileri hariç kaynakla aynı: 4 blok.
- Türkçe gövde aynen korundu; yalnız başa açıklama ve dil bağlantısı eklendi.
- Markdown kod çitleri dengeli; iki dilde AI notları mevcut.
- Göreli belge bağlantılarının kök girdileri GitHub dosya listesinde var; README-TR.md bu commit ile ekleniyor.
- İlk kontrol çağrısı doğrulama ifadesindeki fazla parantez nedeniyle sözdizimi hatası verdi. İfade düzeltildikten sonra tüm kontroller geçti; uygulama hatası değildir.

Sonuç: belge kontrolleri geçti. Araç JavaScript kontrolleri için süreç çıkış kodu yoktur. Uygulama, derleme, kurulum ve gerçek OBS testleri çalıştırılmadı; yalnız dokümantasyon değişti. Aktif yerel kurulum değiştirilmedi. Diğer rehberlerin dili değiştirilmedi; İngilizce README'de ilgili Türkçe rehber bağlantıları belirtildi.
