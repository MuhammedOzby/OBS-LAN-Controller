# Deck tarayıcı incelemesi

Tarih/saat (UTC): 2026-10-29T12:52:41.842089+00:00

Amaç: Düzenleyici hedef kaybı, mobil yerleşim ve kopmada kilitleme kontrolü.
Ortam: Chromium/CDP, http://127.0.0.1:8080/deck. Sunulan deck.js/deck.css dosyaları kaynakla byte düzeyinde aynı. Fetch sayfa betikleri öncesi fixture ile değiştirildi; gerçek OBS komutu gönderilmedi.

Adımlar ve gerçek sonuçlar:
- Fixture: Main/Break sahneleri, Logo görünürlük, Mic susturma, Clip medya. 12 tuş oluştu.
- Break listeden kaldırıldı; sonraki poll beklendi. İkinci tuş düzenleyicisi açıldı: hedef kendiliğinden Main seçildi, Kaydet aktif, uyarı boş. Yalnız etiket değiştirip kaydetmek localStorage içindeki Break hedefini Main yaptı. Bu istenmeyen yeniden atama yeniden üretildi.
- 1280x720, 844x390, 390x844 boyutlarında belge yatay taşmadı.
- 390x844 portrede grid yüksekliği 181.516px iken dört satırın her biri 108px. Son tuş altı 660.922px, workspace altı 450.281px; alt not 398.438px ve sidebar başlangıcı 462.281px. Tuşlar alt içerikle çakışıyor. Yalnız yatay taşma testi bu hatayı yakalamaz.
- connected=false sonrası tüm tuşlar kilitlendi, Kontroller kilitli yazısı görüldü. POST sayısı 0.

Sınırlamalar: fixture simülasyonudur; gerçek OBS medya/sahne/ses işlemleri, fiziksel telefon ve gecikme/yük testi yapılmadı. Görsel piksel incelemesi değil DOM/geometri doğrulaması.
