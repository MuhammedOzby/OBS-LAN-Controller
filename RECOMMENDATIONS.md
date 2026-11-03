# OBS için kolay başlangıç

Önce OBS'nin kendi içinden kısa bir kayıt yapın. OBS kaydedemiyorsa panel de kaydı başlatamaz.

## 1. OBS bağlantısını açın

1. OBS'de **Araçlar → WebSocket Sunucusu Ayarları** bölümünü açın.
2. Sunucuyu ve parola ile doğrulamayı etkinleştirin. Varsayılan port **4455**.
3. Controller'ı açıp bu bilgisayarda **Ayarlar** penceresine aynı port ve parolayı girin. İki uygulama aynı bilgisayardaysa OBS adresi **127.0.0.1** olsun. Dosyayı elle düzenlemeniz gerekmez.

OBS parolasını telefona yazmayın veya paylaşmayın. Telefon erişimi için README'deki eşleştirme kodunu kullanın. Ayrı obs-websocket eklentisi kurmanız gerekmez.

## 2. Kayıt başlamıyorsa

1. OBS'de **Ayarlar → Çıkış → Kayıt** bölümünü açın. Basit çıkış modunda kayıt seçenekleri aynı Çıkış sayfasındadır.
2. **Kodlayıcı/Encoder** alanında bilgisayarınızın desteklediği bir seçenek seçili olsun. Alan boşsa veya seçilen kodlayıcı çalışmıyorsa kayıt başlayamaz. Her bilgisayarda aynı kodlayıcı bulunmaz.
3. Kayıt klasörünün var olduğundan, yazılabilir olduğundan ve diskte boş yer bulunduğundan emin olun.
4. OBS'nin **Kaydı Başlat** düğmesiyle birkaç saniyelik deneme yapın ve durdurun. Dosyanın açıldığını kontrol edin; sonra panelden deneyin.

**İşlem hata verirse:** OBS'deki gerçek kayıt durumunu kontrol edin; aynı düğmeye art arda basmayın. Hata, işlemin OBS'de hiç uygulanmadığı anlamına gelmeyebilir. OBS'yi zorla kapatmak kayıt dosyasını bozabilir.

## 3. Kaydı duraklatamıyorsanız

OBS bazı çıkış ayarlarında duraklatmayı desteklemez. Basit modda **Yayın ile aynı** kayıt kalitesi veya gelişmiş modda yayın kodlayıcısını paylaşan kayıt buna neden olabilir. Duraklatma gerekiyorsa bağımsız bir kayıt kalitesi/kodlayıcısı seçip önce OBS içinde deneyin. Ayar adları OBS sürümüne ve dile göre değişebilir.

## 4. Telefonla bağlanırken

- Telefon ve bilgisayar aynı güvenilir yerel ağda olsun. Ayrıntılı adımlar: [README](README.md#telefonlan-eşleştirme).
- İnternete port açmayın. Mevcut HTTP erişimi şifreli değildir; yalnız güvendiğiniz ağda kullanın. Güvenilir HTTPS ve tam PWA kurulumu ayrı teslimdir.
- Önizleme eklentisi kontrol bağlantısından ayrıdır. Görüntü yokken OBS bağlı olabilir; ilk teslimde video eklentisini kurmak zorunlu değildir.

Genel bir kodlayıcı veya kalite ayarı her bilgisayar için garanti değildir. Önce kısa kayıtla doğrulayın; uygulama sistem veya güvenlik duvarı ayarlarını kendiliğinden değiştirmez.

**Bilinen sınırlama:** hızlı kayıt/duraklat/devam/durdur testinde kayıt durdurma takılması gözlendi; çözülmüş değildir. Ayrıntılar [kayıt test raporunda](test_output/2026-11-01_19-10-18_record-live-fast-failure.md). Bu rehber bir sorunsuz çalışma garantisi değildir.
