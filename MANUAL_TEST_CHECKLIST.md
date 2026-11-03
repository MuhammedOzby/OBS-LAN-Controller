# Yarın için elle test listesi

Amaç: otomatik testleri geçen yeni sürümü gerçek OBS ve telefonla doğrulamak. Her bölümde sonucu **Geçti / Olmadı / Denemedim** diye işaretleyin. Bir sorun çıkarsa o bölümde durun; sonraki denemeler için OBS'nin durumundan emin olun.

## 0. Güvenlik ve önceki kayıt — ilk bunu yapın

- [ ] OBS'de gerçek yayın kapalı. Yayın açıksa testleri erteleyin; sırf test için mevcut yayını durdurmayın.
- [ ] OBS'deki kayıt durumuna baktım. Önceki otomatik testten kalan kayıt aktif/takılı görünebilir.
- [ ] Aktif kaydın önemli bir kullanıcı kaydı olmadığından eminsem OBS içinden **Kaydı Durdur** düğmesine bir kez bastım. Tamamen durduğunu gördüm.

**Beklenen:** OBS'de düğme yeniden **Kaydı Başlat** olmalı; kayıt devam ediyor/durduruluyor görünmemeli. Takılı kalırsa yeni kayıt testi yapmayın, düğmeye art arda basmayın ve OBS'yi zorla kapatmayın. Görünen mesajı not edip bana bildirin. Sahne/ses testlerini de kayıt durumunu netleştirene kadar erteleyin.

## 1. Yeni sürümü açın

- [ ] `OBS LAN Controller` klasöründe terminal açıp `go run .` çalıştırdım.
- [ ] Bilgisayardaki tarayıcıda **http://127.0.0.1:8080** adresini açtım. Farklı HTTP portu ayarlıysa onu kullandım.
- [ ] Panel **OBS bağlı** gösteriyor; OBS FPS/CPU gibi veriler görünüyor.
- [ ] **Kontrol**, **Deck**, **Yayın Yönetimi**, **Mobil cihaz ile giriş** adlı dört sekme var.

**Not:** eski `obs-lan-controller` dosyasına çift tıklamak yerine `go run .` kullanın; böylece yeni kaynak kodu çalışır. Port kullanımda hatası çıkarsa rastgele süreç kapatmayın; mesajı not edin.

## 2. Ayarlar ve parolanın korunması

- [ ] Bilgisayardaki panelde **Ayarlar** düğmesi var ve pencere açılıyor.
- [ ] Kayıtlı OBS parolası açık metin veya doldurulmuş parola alanı olarak görünmüyor.
- [ ] Yeni parola alanını boş, **Kayıtlı parolayı sil** seçeneğini işaretsiz bırakarak kaydettim. Diğer değerleri değiştirmedim.
- [ ] Panel yeniden bağlanıyor ve sayfayı yenileyince de bağlı kalıyor.

**Beklenen:** boş parola alanı mevcut parolayı silmemeli. Bu test için gerçek parolanızı değiştirmeyin veya silmeyin.

## 3. Üç favori sahne

Önce OBS'deki mevcut program sahnesinin adını not edin: __________

- [ ] **Kontrol → Favorileri düzenle** penceresinde üç yuva görüyorum.
- [ ] OBS'de varsa ilk üç dışındaki bir sahneyi de seçip kaydettim. Aynı sahneyi iki yuvaya seçmedim.
- [ ] Sahne düğmesine basınca OBS'nin program sahnesi gerçekten değişti ve panel doğru aktif sahneyi gösterdi.
- [ ] Bir yuvayı **Boş** kaydedince o düğme kullanılamaz oldu.
- [ ] Sayfayı yenileyince seçimler korundu.
- [ ] Controller'ı terminalde `Ctrl+C` ile durdurup tekrar `go run .` çalıştırınca seçimler yine korundu.
- [ ] İşim bitince önceki program sahnesine döndüm.

**Beklenen:** her zaman tam üç sahne düğmesi var; seçilen sahne ile OBS'de açılan sahne aynı. Yeterli sahne yoksa ilk üç dışındaki seçim adımını **Denemedim** yazın; sırf test için mevcut sahneleri silmeyin.

## 4. Ses seviyesi ve susturma

Denenecek kaynak: __________  İlk seviye: __________  Başlangıçta susturulmuş mu: __________

- [ ] Kontrol ekranında OBS'nin ses kaynaklarını görüyorum.
- [ ] Bir kaynağın seviyesini örneğin %30'a getirdim; OBS ses mikserindeki değeri de değişti. OBS değeri dB cinsinden gösterebilir; yüzdeyle aynı yazması gerekmez.
- [ ] Kaydırıcıyı birkaç kez sürükleyip bıraktım; son bıraktığım değer geri sıçramadı.
- [ ] Panelden susturdum ve sesi tekrar açtım; OBS'deki susturma durumu aynı oldu.
- [ ] OBS içinden seviyeyi/susturmayı değiştirdim; birkaç saniye içinde panele yansıdı.
- [ ] Kaynağı başlangıçtaki seviye ve susturma durumuna geri getirdim.

**Beklenen:** değişiklik iki tarafta tutarlı. Kaynak görünmüyorsa veya ses destekli kaynak yoksa not edin; rastgele mikrofon/cihaz kurulumu yapmayın.

## 5. Bağlantı kopması ve geri gelmesi

Bunu kayıt/yayın kapalıyken yapın.

- [ ] OBS'de **Araçlar → WebSocket Sunucusu Ayarları** bölümünde sunucuyu geçici kapattım ve değişikliği uyguladım.
- [ ] Birkaç saniye sonra panel bağlantı kaybını gösterdi; sahne/ses ve kritik işlem düğmeleri kullanılamaz oldu. Eski metrikler güncelmiş gibi kalmadı.
- [ ] OBS WebSocket sunucusunu aynı ayarlarla tekrar açtım.
- [ ] Panel kendiliğinden bağlandı ve gerçek durumu yeniden gösterdi.
- [ ] Yeniden bağlanınca kendiliğinden kayıt/yayın başlamadı veya eski bir işlem tekrar uygulanmadı.

**Beklenen:** kopmada güvenli pasif durum, geri gelince doğru veri. OBS'nin kendisini kapatmak gerekmez.

## 6. Telefon ve eşleştirme

- [ ] Telefon ve bilgisayar aynı güvenilir Wi-Fi/LAN'da.
- [ ] HTTP dinleme adresi gerekiyorsa bilgisayardaki **Ayarlar** penceresinde `0.0.0.0` olarak kaydedildi ve yalnız Controller yeniden başlatıldı. Portu değiştirmedim.
- [ ] Bilgisayarda **Mobil cihaz ile giriş** sekmesini açtım; QR ve elle giriş kodu görünüyor. Aynı seçenekler **http://127.0.0.1:8080/pair** sayfasında da var.
- [ ] Telefonda normal kamera ile kendi Wi-Fi/LAN adresimin altındaki QR'ı okutup bağlantıyı açtım. Kod yazmadan Kontrol ekranına geçtim.
- [ ] İkinci bir denemede yeni kod aldım ve telefonda **http://BİLGİSAYARIN_LAN_IP_ADRESİ:8080** üzerinden elle kodla giriş yedeğini denedim; `0.0.0.0` adresini kullanmadım.
- [ ] Telefonda dört sekme, sahne düğmeleri ve ses paneli kullanılabilir. Dikey/yatay çevirdiğimde düğmeler üst üste binmiyor veya erişilemez olmuyor.
- [ ] Telefonda Mobil cihaz ile giriş sekmesinde yeni QR/kod veya oturum iptal düğmesi gösterilmiyor.
- [ ] Telefonda OBS bağlantı/parola **Ayarlar** düğmesi görünmüyor.
- [ ] Bilgisayardaki `/pair` sayfasından **Tüm LAN oturumlarını iptal et** yaptım. Telefonda sayfayı yenileyince tekrar kod istendi.

**Beklenen:** eşleşmemiş telefon kontrol edemez; eşleşmiş telefon da yönetim/parola ayarına erişemez. İnternete port açmayın, güvenlik duvarını tümden kapatmayın. HTTP şifreli değildir. Görüntü önizlemesi bu listenin zorunlu kabulü değildir.

**QR notu:** kod 5 dakika geçerli ve tek kullanımlık. Eski QR ile giriş reddedilirse bilgisayarda yeni kod açın. Controller yeniden başlatılınca telefon oturumu silinir; bu durumda yeniden QR okutmak gerekir. Sonuç gönderirken QR'ın veya kodun göründüğü ekran görüntüsünü paylaşmayın.

## 7. Kısa kayıt — yalnız 0. bölüm sorunsuzsa

OBS'de başka kayıt ve yayın olmadığından emin olun. Önemli içerikle denemeyin. Düğmeler doğrudan çalışır; ek onay penceresi yoktur.

- [ ] Önce OBS'nin kendi düğmesiyle yaklaşık 10 saniyelik kayıt yaptım, durdurdum ve dosyanın açıldığını kontrol ettim.
- [ ] Panelde **Yayın Yönetimi → Kaydı Başlat** kullandım; OBS'de başladığını gördüm ve en az 5 saniye bekledim.
- [ ] Duraklatma destekleniyorsa duraklattım, OBS'de durakladığını gördüm ve birkaç saniye bekledim.
- [ ] Devam ettirdim, OBS'de devam ettiğini gördüm ve en az 5 saniye bekledim.
- [ ] **Kaydı Durdur** kullandım; hem OBS hem panel kaydın tamamen kapandığını gösterdi.
- [ ] Oluşan dosyayı açıp oynatabildim.

**Beklenen:** panel ve OBS durumları aynı. Duraklatma desteklenmiyorsa bu adımı **Desteklenmiyor** diye not edin; sırf test için encoder/profil değiştirmeyin. Hata veya **Durduruluyor** takılması olursa tekrar başlatmayın ve zorla kapatmayın. Görünen durumları not edin. Hızlı art arda başlat/duraklat/devam/durdur denemesi şimdilik yapmayın.

**Gerçek yayın başlatma/durdurma bu listeye dahil değil.** Platformunuza sırf test için canlı yayın açmayın; bunu daha sonra güvenli bir hedefle ayrıca doğrularız.

## Bana göndereceğiniz kısa sonuç

Parola, eşleştirme kodu, çerez veya yayın anahtarı göndermeyin. Ekran görüntüsünde bunlar varsa gizleyin.

```text
Bilgisayar tarayıcısı:
Telefon modeli / tarayıcısı:
0. Önceki kayıt: Geçti / Takılı (OBS'deki tam yazı: ...)
1. Açılış ve bağlantı:
2. Ayarlar/parola korunması:
3. Favori sahneler:
4. Ses:
5. Kopma/yeniden bağlanma:
6. Telefon/eşleştirme:
7. Kayıt ve dosyayı oynatma:
Sorun çıkan adım ve ekrandaki tam mesaj:
```
