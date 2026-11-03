# Page heading — tarayıcı ve çalışan sunucu

Tarih: 2026-10-29T14:18:02.221977+00:00

## Tarayıcı
127.0.0.1:18080 test sunucusu; OBS bağlantısı yok.
Üç sayfa × 1280×720, 844×390, 390×844: 9 DOM kontrolü geçti. Başlık/etiket yok, yatay taşma yok; içerik main alanını dolduruyor.
İlk deneme çıkış 1: Deck gizli dialogu yanlışlıkla yerleşim çocuğu sayıldı; test in-flow çocukları sayacak şekilde düzeltildi. İkinci deneme çıkış 1: rapor yazarken tanımsız workspace değişkeni; uygulama hatası değil. Üçüncü deneme çıkış 0; 9 kontrol geçti.

## Çalışan uygulama
Yeni derlenen test_output/obs-lan-controller-final ile mevcut 8080 sunucusu yeniden başlatıldı.
/: HTTP 200, başlık bloğu yok
/deck: HTTP 200, başlık bloğu yok
/yayin-yonetimi: HTTP 200, başlık bloğu yok
8080 üzerinde sunulan CSS kaynak dosyayla aynı.
HTTP readback çıkış kodu: 0.
OBS kontrol komutu gönderilmedi; gerçek OBS davranışı bu değişiklik kapsamında denenmedi.
