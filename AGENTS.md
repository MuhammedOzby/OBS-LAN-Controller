# OBS LAN Controller — ajan talimatları

## Kimlik ve sınır

Repo adı: `obs-lan-controller`. Bu repo Go CLI ve tarayıcı panelidir. OBS eklentisi ayrı repodadır; onu burada uygulamayın. Yalnızca bu klasörde çalışın. Kullanıcı izin vermedikçe üst veya kardeş klasörlere erişmeyin. Ortak talimatlar görev bağlamında verildiyse onlara da uyun.

Önce `PLAN.md` dosyasını okuyun. İlk iş Milestone 1'dir; diğer aşamalara kendiliğinden kapsam genişletmeyin.

## Kod işlerinde delegasyon

- Kullanıcının 2 Kasım 2026 tarihli güncel tercihiyle kod işleri **GPT-6.1 Terra** (`openai/gpt-6.1-terra`) alt ajanına devredilir. Önceki Ollama/yalnız yerel model zorunluluğu bu proje için kaldırılmıştır; Ollama servisini onarmak bu görevin parçası değildir.
- Ana koordinatör inceleme, test çalıştırma, planlama ve belge güncellemesi yapar; uygulama/test kodunu ve düzeltmelerini Terra alt ajanına yaptırır. Alt ajan yönetimi koordinatöre bırakılmıştır.
- Bu kod yazma sınırı yalnız ana koordinatöre aittir. Görevi devralan Terra alt ajanı uygulama ve test kodunu kendi araçlarıyla yazar; yeniden alt ajan başlatması gerekmez.
- Her alt ajan bu repo sınırları, gizli bilgi koruması ve `test_output/` raporlama kurallarına uyar. Bu tercih OBS eklentisi reposuna erişim vermez.

## Teknik tercihler

- Backend Go. Kullanıcının tercihi doğrultusunda web yönlendirmesi ve middleware için Gin Web Framework, sayfalarda güvenli `html/template` kullanın; şablon uzantısı `.tmpl` olsun. Zaman aşımları ve temiz kapanma için Gin'i `http.Server` üzerinde çalıştırın.
- HTML, CSS, JavaScript ve PWA dosyalarını `go:embed` ile dağıtılabilir uygulamaya gömmeyi hedefleyin.
- Tarayıcıda anlık durum, ses kontrolü ve WebRTC için gerekli JavaScript bulunabilir. React/SPA veya Node.js çalışma zamanı zorunluluğu eklemeyin.
- Linux ilk hedeftir. İşletim sistemine özgü dosya yollarını, servis yönetimini ve izin işlemlerini ayrı tutun.
- Go sunucusu OBS WebSocket istemcisidir. Güncel OBS'nin dahili sunucusunu kullanın; varsayılan port 4455, yerel adres 127.0.0.1 olsun. Kullanıcı bunları değiştirebilsin.
- WebRTC için Pion adaydır. Video yakalama/kodlama sunduğunu varsaymayın; kodlanmış videoyu eklentiden alıp taşıyın. Controller'da yeniden kodlamayı varsayılan hale getirmeyin.
- Önerilen CLI komutları: `setup`, `serve`, `doctor`. İlk teslimde yalnızca gerekli komutları gerçekleştirin; olmayan komutu çalışıyormuş gibi göstermeyin.

## Arayüz kuralları

Arayüz teması olarak **Bootstrap Kehribar Gece** kullanılmalıdır: https://github.com/MuhammedOzby/bootstrap-kehribar-gece. `.tmpl` tabanlı paneli bu temayla geliştirin; kullanıcı onayı olmadan başka bir temayla değiştirmeyin.

1. **Kontrol:** Program görüntüsü, performans metrikleri, tam üç favori sahne yuvası, ses kontrolleri. Başlat/durdur yok.
2. **Deck:** Düzenlenebilir tuş ızgarası. Sahne, kaynak görünürlüğü, susturma ve medya işlemleri. Yayın/kayıt başlatma-durdurma veya bunu dolaylı yapan genel komut/makro atanamaz.
3. **Yayın Yönetimi:** Yayın/kayıt başlatma-durdurma, destekleniyorsa kayıt duraklatma/devam, durum ve süreler. Kullanıcı tercihiyle düğmeler ek onay ekranı olmadan doğrudan çalışır; işlem kilidi ve gerçek OBS durum doğrulaması korunur.
4. **Mobil cihaz ile giriş:** Kullanıcının 2 Kasım 2026 tarihli isteğiyle eklenen eşleştirme sekmesi. OBS bilgisayarında kısa ömürlü tek kullanımlık QR/kod gösterir; LAN istemcisi yeni eşleştirme yetkisi üretemez veya okuyamaz. Önceki tam üç sekme kuralı bu istekle dört sekme olarak değiştirilmiştir.

- Bağlantı ayarları ortak ayarlar penceresinde kalır; dördüncü sekme yalnız mobil eşleştirme içindir.
- Gerçek durumu OBS'den alın. Komutu göndermek başarı sayılmaz; yanıt ve durum değişimini izleyin.
- Bağlantı kopunca kontrolleri pasifleştirin. Yeniden bağlanınca durumu tekrar okuyun; eski kritik komutları yeniden oynatmayın.
- Yayın/kayıt işlemlerinde kör toggle yerine açık başlat/durdur çağrıları kullanın. Tekrarlı dokunmaları engelleyin.
- Ses kaydırıcısı isteklerini sınırlayın; dokunmatik kullanım ve erişilebilir adlandırma sağlayın.
- OBS FPS ile önizleme FPS'sini ayrı gösterin. GPU kullanımı gibi standart API'nin sağlamadığı metrikleri uydurmayın.

## Bağlantı ve gizli bilgiler

- OBS parolasını sunucu tarafında tutun; tarayıcıya, HTML'e, loglara veya Git'e taşımayın. Parolayı CLI argümanıyla vermeyi zorunlu tutmayın.
- İlk kurulum bilgisayarda/loopback üzerinden yapılmalı; LAN'dan gelen eşleştirilmemiş kullanıcı OBS parolasını veya yönetim ayarlarını görememelidir.
- Telefon erişimi için kısa ömürlü eşleştirme ve doğrulanmış oturum tasarlayın. QR kodda kalıcı OBS parolası bulunmasın.
- `/pair` ve mobil giriş sekmesi aynı güvenli eşleştirme akışını kullanır. QR internet servisine gönderilmeden yerelde üretilir; LAN'dan erişilebilir gerçek sunucu adresini içerir. Kod erişim günlüklerine/Referer'a taşınmamalı, süre sonu ve tek kullanımlılık korunmalıdır. Tarayıcı kamerası/HTTPS zorunluluğu eklemek yerine telefonun normal kamerasıyla okutma desteklenir.
- Yerel tarayıcı kurulumu ile PWA kurulumu farklıdır: LAN HTTP çalışıyor diye PWA kurulumunun tamamlandığını iddia etmeyin. Güvenilir HTTPS çözümü üretim teslim şartıdır.
- Bağlantı türünü sayfanın HTTPS olmasından türetip OBS adresini otomatik WSS'ye çevirmeyin; OBS sunucu bağlantısı kendi yapılandırmasına göre kurulmalıdır.
- Oturum, Origin ve CSRF kontrollerini kullanılan HTTP/WebSocket işlemlerine uygun uygulayın. İnternete port açılmasını gerektirmeyin.

## Doğrulama ve teslim

Durum makinesi, yeniden bağlantı, kritik komut koruması, metrik hesapları ve eşleştirme gibi davranışları anlamlı testlerle doğrulayın. Arayüzü mobil ve masaüstünde kontrol edin. Gerçek OBS yoksa simüle edilen ve doğrulanamayan davranışları belirtin.

Proje boyunca her test veya test grubu çalıştırmasının sonucunu repo kökündeki `test_output/` klasöründe bir Markdown raporuna kaydedin. Başarısız ve engellenen denemeleri de raporlayın; önceki sonuçları silmeyin. Her raporda tarih/saat, amaç, çalıştırılan komutlar veya manuel adımlar, gerçek sonuç ve çıkış kodu (varsa), simülasyon/gerçek OBS ayrımı ve doğrulanamayan noktalar bulunsun. Gizli bilgileri raporlara yazmayın. Test çalıştırılmadıysa bunu açıkça belirtin; depo/belge kontrollerini uygulama testi gibi sunmayın. Ayrıntılar için `TESTING.md` dosyasını okuyun.

`PLAN.md` ilerlemesini güncelleyin. Son raporda değişiklikleri, doğrulamayı ve açık sınırlamaları kısa yazın.
