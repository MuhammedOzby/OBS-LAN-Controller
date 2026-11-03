# Milestone 1 kabul kontrolü

Bu liste bir doğrulama planıdır; işaretlenmeyen maddeler tamamlandı sayılmaz. Testlerin tarihli sonuçları `test_output/` altında tutulur. Yalnız bu teslimin kapsamı değerlendirilir; mevcut Deck/video çalışmalarının varlığı bunların yeni kabulünü kanıtlamaz.

## Güvenli test sırası

1. Otomatik Go/JavaScript testleri ve `go vet ./...`.
2. Yerel başlangıç ve ayarlar: eksik dosya, yanlış adres, yanlış parola, başarılı bağlantı. Parola yanıt, şablon veya raporda görünmemeli.
3. LAN: eşleştirilmemiş erişim engeli, QR ile otomatik ve elle kodla erişim, süre sonu/tekrar kullanım, oturum iptali; LAN'dan kod/QR üretimi ve yönetim ayarları engeli. QR'ın doğru LAN adresi/portuna açılması ve gerçek telefon kamerasıyla okutma ayrıca kontrol edilir.
4. Sahne: tam üç favori seçimi, yeniden açınca korunması, aktif sahne, eksik/silinmiş sahne ve OBS dışından değişiklik.
5. Ses: geçici ses girdisinde seviye/susturma, komut sonrası tekrar okuma, OBS dışından değişikliğin panele yansıması. Gerçek kullanıcı girdisinin ayarları korunur.
6. Kayıt: mevcut kayıt yoksa kısa başlat/duraklat/devam/durdur döngüsü. Desteklenmeyen duraklatma açıkça anlatılmalı. Zaten devam eden kullanıcı kaydı durdurulmaz.
7. Yayın: mevcut halka açık yayın başlatılmaz/durdurulmaz. Canlı test için izole yerel alıcı kurulabiliyorsa kullanılır ve geçici ayarlar geri alınır; aksi halde simülasyon ile gerçek OBS sınırı raporlanır.
8. Bağlantı: kopunca eski veriler/komutlar kullanılmaz, tekrar bağlantıda durum yeniden okunur; kritik komutlar tekrar oynatılmaz.
9. Tarayıcı: kullanıcının güncel tercihiyle dört sekme; dördüncüsü Mobil cihaz ile giriş. Yalnız Yayın Yönetimi'nde yayın/kayıt işlemleri, mobil/tablet/masaüstü yerleşim ve erişilebilir etiketler.

## Teslim sınırları

- Fiziksel telefon/tablet testi, tarayıcı boyutu simülasyonuyla aynı değildir.
- Video eklentisi, tam PWA ve üretim HTTPS bu aşamanın teslimi değildir. HTTP'nin şifreli olmadığı kullanıcıya anlatılmalıdır.
- Ekran görüntüsü döngüsü nihai video diye sunulmaz. Mevcut önizleme korunur; eklenti deposuna bu görevde girilmez.
- OBS veya uygulama ayarı değiştirilirse neden ve basit kullanıcı adımları `RECOMMENDATIONS.md`/README'ye yazılır. Test kaydının yeri kullanıcıya bildirilir; mevcut dosyalar silinmez.
- PLAN kutuları yalnız karşılık gelen uygulama ve doğrulama gerçekten bittiğinde işaretlenir.
