# Manuel kabul denemesi — 2 Kasım 2026

Saat dilimi: Europe/Istanbul (UTC+03:00). Başlangıç yaklaşık 14:50; rapor 15:05.

Amaç: MANUAL_TEST_CHECKLIST.md adımlarını mevcut gerçek OBS ve Chrome ile uygulamak.

Ortam: Kullanıcının mevcut Chrome tarayıcısı; localhost:8080 yönetim paneli ve gerçek LAN adresinde ayrı istemci sekmesi. Panel gerçek OBS'ye bağlıydı. OBS'nin yerel penceresini kontrol eden araç bu oturumda mevcut değildi. Fiziksel telefon yoktu. Simüle OBS kullanılmadı; sonuçlar panelin OBS'den okuduğu durumlarla sınırlıdır. Uygulama/test kodu değiştirilmedi; otomatik test çalıştırılmadı.

## Sonuçlar

| Bölüm | Sonuç | Kanıt / açık kalan |
|---|---|---|
| 0. Güvenlik | Kısmen geçti | Yayın Yönetimi: OBS bağlı, YAYIN KAPALI, KAYIT KAPALI, Hazır; Kaydı başlat etkin, durdur/duraklat/devam pasif. Aktif kayda müdahale edilmedi. OBS'nin kendi düğmesi doğrudan görülmedi. |
| 1. Açılış | Kısmen geçti | Mevcut sunucu açıldı; dört sekme ve güncel metrikler görüldü (OBS FPS 30.00, değişen CPU/render). Bu oturumda go run . başlatılmadı; çalışan sürecin kaynak sürümü doğrulanamadı. Port çatışması denenmedi. |
| 2. Ayarlar/parola | Geçti (panel) | Ayarlar açıldı. Yeni parola boş; Kayıtlı parolayı sil işaretsiz; Kayıtlı parola var bildirimi mevcut. Hiçbir alan değiştirilmeden Kaydet: Kaydedildi. Dinleme adresi değişikliği için sunucuyu yeniden başlatın. Sonrasında kapat/yenile: OBS bağlı, metrikler canlı. Kayıtlı parola okunmadı/değiştirilmedi. |
| 3. Favoriler | Kısmen geçti | Başlangıç programı Scene; favoriler Ara ekran / Scene / Konuşma. Tam üç yuva görüldü. Yuva 1 Boş kaydedildi; Sahne 1 yok pasif, yenilemede korundu. Ara ekran tekrar kaydedildi. Ara ekran seçilip yenilemede program ve aktif işaret Ara ekran oldu. Sonrasında Scene geri seçildi ve gerçek durum tekrar Scene olarak okundu. Yalnız üç sahne olduğundan ilk üç dışındaki sahne denenmedi. Süreç yeniden başlatma ve OBS penceresinden bağımsız görsel doğrulama yapılmadı. |
| 4. Ses | Kısmen geçti | Desktop Audio başlangıç %100, susturulmamış. %30 ayarlandı, yenilemede korundu. 45→65→30 ardışık değerleri girildi; son değer korundu. Susturma etkinleştirildi, Ses durumu doğrulandı ve yenilemede korundu. Bu işlem fareyle sürükleme değil erişilebilir kaydırıcı değer değişimidir. OBS mikseri ve OBS içinden ters yön değişiklik denenmedi. Temizlikte %100 ayarını hemen izleyen ilk susturma tıklaması durumu değiştirmedi; son kontrolde susturma hâlâ açıktı. Ayrı tıklamayla kapatıldı, yenilemede %100 ve susturulmamış durum doğrulandı. Bu ara tıklamanın nedeni araştırılmadı. |
| 5. Kopma / yeniden bağlanma | Denemedim | OBS WebSocket sunucusunu OBS içinden kapatıp açmak için yerel uygulama erişimi yok. OBS/süreç zorla kapatılmadı; sunucu ayarları değiştirilmedi. |
| 6. Telefon / eşleştirme | Kısmen geçti; tamamlanamadı | LAN adresine eşleşmemiş Chrome sekmesi /pair kod formuna yönlendi; kontrol paneli açılmadı. Yerel mobil girişte iki LAN adresi için QR, süre ve kod görüldü. Panelin kendi bağlantısıyla yerel /pair açıldı. Kod, rapora/loga eklenmeden yerel DOM'dan LAN giriş alanına aktarıldı; Eşleştir sonrası Kontrol açıldı, OBS bağlı ve dört sekme vardı. LAN'da Ayarlar yok; Mobil cihaz ile giriş sekmesinde yeni QR/kod ve iptal düğmesi yok, yalnız giriş formu var. Fiziksel kamera/telefon, QR ile otomatik giriş, yeniden kullanım ve iptal testleri yapılmadı. |
| 7. Kısa kayıt | Denemedim | Önce OBS'nin kendi düğmesiyle 10 saniyelik kayıt ve dosya oynatma önkoşulu uygulanamadı. Sonradan tarayıcı erişimi de kesildi. Kayıt/yayın komutu gönderilmedi; dosya oluşturulmadı. |

## Gözlemler ve engeller

1. LAN sekmesi açıldıktan sonra yerel /pair doğrudan gezinme ile açıldığında Forbidden Origin görüldü; yenileme aynı sonucu verdi. Yerel Kontrol → Mobil cihaz ile giriş → eşleştirme sayfasını yenileyin bağlantısı ise /pair sayfasını başarıyla açtı. Güvenlik kontrolü devre dışı bırakılmadı. Doğrudan gezinmenin istek başlıkları ölçülmediği için bunu kesin uygulama hatası olarak sınıflandırmıyoruz; gezinme yoluna bağlı gözlem.
2. Mobil ekran denemesi için 390×844 viewport isteğiyle aynı çağrıda Kontrol'e gezinme ve durum/screenshot istendi. Araç 30 saniyede zaman aşımına uğradı ve oturumu sıfırladı; hangi alt adımın tamamlandığı doğrulanamadı. Ardından mevcut LAN sekmesine yeniden bağlanma da 30 saniyede zaman aşımına uğradı. Yerleşim başarılı sayılmadı; viewport değişikliği uygulanmış olabilir ve geri alındığı doğrulanamadı.
3. Son başarılı sunucu okumaları: Scene aktif; favoriler Ara ekran / Scene / Konuşma; Desktop Audio %100, susturulmamış; diğer ses kaynaklarına dokunulmadı. Yayın/kayıt başlatılmadı. Yeni test LAN oturumunun iptali araç kesintisi nedeniyle yapılamadı; uygulama kuralına göre 12 saat sonra veya Controller yeniden başlatılınca biter. Kullanıcının diğer LAN oturumlarına müdahale edilmedi.
4. Yerel terminal ortamı, kullanıcının masaüstü süreçlerini listelemiyordu; bu nedenle Controller yeniden başlatma denenmedi. Yanlış sürece sinyal gönderilmedi.

Manuel işlemlerin işletim sistemi çıkış kodu yoktur. Belge okumaları başarılıydı; tanılama dosya aramasında mevcut olmayan Lib/security.go için çıkış 2 alındı, doğru dosya Lib/http_security.go olarak listelendi. Bu, uygulama testi hatası değildir. İki eşleştirme otomasyonu hazırlama denemesi DOM metin eşleşmesi/değişken kapsamı nedeniyle form gönderilmeden durdu; üçüncü deneme gerçek LAN girişini başarıyla tamamladı.

## Kısa sonuç

Bilgisayar tarayıcısı: Kullanıcının mevcut Chrome'u; sürüm ölçülmedi.
Telefon modeli / tarayıcısı: Fiziksel telefon kullanılmadı; ayrı Chrome LAN sekmesi.
0. Önceki kayıt: Panelde kapalı / hazır; OBS penceresi kontrol edilmedi.
1. Açılış ve bağlantı: Bağlantı geçti; yeni sürümü başlatma denenmedi.
2. Ayarlar/parola korunması: Geçti (panel).
3. Favori sahneler: Panel komutu, boş yuva, yenileme ve geri dönüş geçti; restart açık.
4. Ses: Panel seviye/susturma ve geri yükleme geçti; OBS içinden ters yön ve sürükleme açık.
5. Kopma/yeniden bağlanma: Denemedim.
6. Telefon/eşleştirme: LAN kod girişi ve yönetim kısıtları geçti; fiziksel cihaz, QR girişi, yerleşim ve iptal açık.
7. Kayıt ve dosyayı oynatma: Denemedim.

Tam manuel kabul verilmedi. Kalan maddeler, OBS yerel pencere erişimi, çalışan tarayıcı aracı ve fiziksel telefonla tamamlanmalıdır.

Son temizlik denemesi: tarayıcı envanteri yeniden okunabildi, ancak ekran boyutu sıfırlama girişimi sekmenin kontrol bağlantısı bağlı olmadığı hatasıyla başarısız oldu. İptal ve ekran boyutu geri alma doğrulanamadı.
