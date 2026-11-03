# OBS LAN Controller

Go ile sunulan, Bootstrap Kehribar Gece tabanlı OBS paneli. **Kontrol** ekranında program önizlemesi, OBS metrikleri, üç favori sahne ve ses seviyesi/susturma bulunur. **Deck** düzenlenebilir 12 tuşla sahne, kaynak görünürlüğü, susturma ve medya kontrolü sunar. **Yayın Yönetimi** yayın/kayıt işlemlerinin tek yeridir. Dördüncü sekme **Mobil cihaz ile giriş**, telefonla QR eşleştirme içindir. Güvenilir HTTPS/PWA henüz yoktur.

## Çalıştırma

OBS'ye yeni başladıysanız önce [kolay OBS ayarları ve kayıt sorunları](RECOMMENDATIONS.md) sayfasına bakın.

Gerçek OBS ve telefonla denemek için: [adım adım elle test listesi](MANUAL_TEST_CHECKLIST.md).

Go 1.24 veya üzeri gerekir. OBS'yi açın; bu klasörde:

```sh
go run .
```

1. Bu bilgisayarda **http://127.0.0.1:8080** adresini açın.
2. İlk çalıştırmada açılan **Bağlantı ayarları** penceresine OBS'nin WebSocket adresini, portunu ve parolasını girip kaydedin. Aynı bilgisayardaysa OBS IP'si **127.0.0.1**, varsayılan port **4455** olur.
3. Telefon erişimi istiyorsanız aynı pencerede **HTTP dinleme IP adresi** alanını **0.0.0.0** yapın. Kaydedin; terminalde `Ctrl+C` ile durdurup `go run .` ile tekrar açın. Bu, güvenilir ağdaki diğer cihazların bağlanmasını sağlar.
4. Telefonda **http://SUNUCU_LAN_IP:8080** adresini açıp aşağıdaki eşleştirme adımlarını izleyin. Bilgisayarın LAN IP'sini `hostname -I` ile öğrenebilirsiniz. **0.0.0.0 tarayıcıya yazılacak adres değildir.**

Dosya yoksa ilk kurulum yalnız bu bilgisayara açık başlar. Dosya zaten varsa kayıtlı adres/port kullanılır. **Ayarlar** düğmesi yalnız OBS bilgisayarından (loopback) görünür; telefondan bağlantı/parola ayarı değiştirilemez. Başka port seçtiyseniz yukarıdaki adreslerde 8080 yerine onu kullanın.

İsteğe bağlı açık komut: `go run . serve`. Ayarlar varsayılan olarak çalışma dizinindeki `configuration.json` dosyasından okunur; başka dizinden çalıştırırken `-config /tam/yol/configuration.json` kullanın.

### Yapılandırma

- `obs.ip`: OBS bilgisayarının IP adresi; başlangıç değeri `127.0.0.1`.
- `obs.port`: OBS WebSocket portu; başlangıç değeri `4455`.
- `obs.password`: OBS WebSocket parolası; kendi parolanızı yerel dosyaya girin. Boş değer parolasız OBS için kabul edilir.
- `server.listening_ip`: HTTP dinleme IP'si; `0.0.0.0` tüm IPv4 arayüzleri, `127.0.0.1` yalnızca bu bilgisayar.
- `server.port`: HTTP portu; başlangıç değeri `8080`.

Normal kullanımda JSON'u elle düzenlemek gerekmez. İsterseniz `configuration.example.json` dosyasını `configuration.json` olarak kopyalayıp düzenleyebilirsiniz. IP adresleri ve 1–65535 port aralığı doğrulanır; mevcut dosya hatalıysa uygulama durur, üzerine yazılmaz. Eksik dosya yerel ilk kurulumu açar. HTTP dinleme adresi/portu değişince yeniden başlatın; pencereden kaydedilen OBS bağlantı ayarı sonraki istekte yeniden bağlanır. `-addr` geçici HTTP adresi önceliğidir; ilk kurulumda yalnız loopback adresi kabul edilir.

`configuration.json` parola içerdiği için Git dışında tutulur; paylaşılabilir şablon `configuration.example.json` dosyasıdır. Panel dosyayı atomik olarak yalnız kullanıcıya açık `0600` izinle yazar. Elle oluşturduysanız `chmod 600 configuration.json` ile koruyun. Mevcut parola tarayıcıya geri gönderilmez; pencerede yeni parola boş bırakılırsa korunur, **Kayıtlı parolayı sil** işaretlenirse temizlenir. OBS'nin yerleşik WebSocket sunucusu etkin olmalıdır. Kopmada veriler temizlenir, sonraki isteklerde iki saniyelik bekleme ile yeniden denenir. Komutlar otomatik tekrar edilmez.

### Sahne ve ses

- Kontrol sekmesinde **Favorileri düzenle** ile üç sahne yuvasını seçip kaydedin. Yuva boş kalabilir; aynı sahne iki kez seçilemez. İlk kayıttan önce ilk üç OBS sahnesi kullanılır.
- Tercihler sunucu tarafında, yapılandırma dosyasının yanındaki `configuration.json.favorites.json` dosyasında saklanır ve eşleşmiş cihazlarca paylaşılır. Silinmiş sahne veya başka sahne koleksiyonu sessizce başka sahneye atanmaz; gerekirse favorileri tekrar düzenleyin.
- Ses bölümünde OBS'nin ses destekli kaynakları listelenir. Kaydırıcı %0–100 arasıdır; düğme susturur/sesi açar. OBS içinde yapılan değişiklikler periyodik okumayla panele yansır. Bağlantı yokken kontroller kapanır; gönderilen işlem OBS'den tekrar okunarak doğrulanır.

Önizleme ayrı OBS eklentisinin H.264 program videosunu Unix soketinden alıp **yeniden kodlamadan WebRTC** ile tarayıcıya taşır. Panel ekran görüntüsü döngüsü kullanmaz. Video bağlantısı ve OBS kontrol bağlantısı ayrı izlenir; video kopması sahne/Deck/yayın kontrollerini kapatmaz. Bit hızı gönderilen baytların zaman farkından hesaplanır; ilk örnek, sayaç sıfırlanması ve pasif/yeniden bağlanan yayın için değer gösterilmez. Kalite etiketi ölçülen düşen kare/ağ tıkanıklığına bağlı teknik göstergedir, izleyici kalitesi ölçümü değildir.

### OBS eklentisiyle aynı bilgisayarda başlatma (Linux)

1. Yerel Linux OBS'ye **OBS LAN Controller Live Preview Extension** eklentisini kurup OBS'yi yeniden açın. Flatpak bu teslimde doğrulanmadı.
2. Controller'ı **OBS'yi açan aynı kullanıcıyla**, masaüstü oturumunun terminalinden çalıştırın. `sudo` kullanmayın. `XDG_RUNTIME_DIR` masaüstü oturumundan gelmelidir (genellikle `/run/user/UID`); rastgele ortak/geçici bir dizin atamayın.
3. Arch/CachyOS'ta Controller derlemesi için `sudo pacman -Syu --needed go` yeterlidir. OBS SDK/CMake/Node/FFmpeg Controller'ın çalışması için gerekli değildir. Bu klasörde:

   ```sh
   # Yalnız yapılandırma dosyanız yoksa:
   cp -n configuration.example.json configuration.json
   chmod 600 configuration.json
   # Dosyadaki obs.ip=127.0.0.1, port ve OBS WebSocket parolasını ayarlayın.
   go build -o obs-lan-controller .
   ./obs-lan-controller serve
   ```

4. Önce `http://127.0.0.1:8080` adresindeki **Kontrol** sekmesini açın. Video, sekme görünür ve WebRTC bağlıyken başlar. Son izleyici ayrılınca veya sayfa gizlenince durur; görünür olduğunda yeniden bağlanır. Eklenti otomatik yeniden bağlandığından Controller'ı OBS'den sonra açabilirsiniz.
5. Telefon için mevcut `/pair` akışıyla eşleştirin ve LAN IP'sinden paneli açın. Modern H.264/WebRTC destekli tarayıcı gerekir. Otomatik oynatma engellenirse video alanına dokunun. Uygulama internet STUN/TURN kullanmaz; yalnız doğrudan yerel/LAN ICE adayları vardır. Güvenlik duvarı HTTP TCP 8080 dışında WebRTC'nin dinamik UDP trafiğine de izin vermelidir. Uygulama güvenlik duvarını değiştirmez, internete port yönlendirmesi yapmayın.

Video soketi IP/port değildir: `$XDG_RUNTIME_DIR/obs-lan-controller/preview.sock`. Controller sunucu, eklenti istemcidir; 0700 dizin/0600 soket ve aynı UID peer kontrolü kullanılır. İkinci Controller bu soketi devralmaz. Soket başlangıcı başarısızsa HTTP/OBS kontrol çalışmaya devam eder ve logda neden gösterilir. Durum teşhisi için yerel tarayıcıda `/api/preview/status` açılabilir. `connected` eklenti bağlantısını, `streaming` son üç saniyede gelen H.264 video verisini belirtir; tarayıcının gerçekten çözdüğü kare veya gecikme ölçümü değildir. API'deki `fps` encoder yapılandırmasıdır; paneldeki önizleme FPS'i ise tarayıcı WebRTC alıcı istatistiğinden ölçülür ve veri yoksa `—` gösterilir.

İlk sürüm en fazla **3 eşzamanlı önizleme oturumu** kabul eder; bu bir ölçülmüş kapasite garantisi değil kaynak sınırıdır. Bağlanmamış oturumlar 15 saniyede, oturumlar en geç bir saatte kapatılır; görünür panel yeniden bağlanır. LAN eşleştirme iptali/süre bitişi video oturumunu da kapatır. Yavaş izleyicide sınırlı kuyruk temizlenip sonraki anahtar kare beklenir. Ses gönderilmez. Mevcut eklenti protokolü değişmediğinden bu Controller güncellemesi için eklentiyi yeniden kurmak gerekmez.

Doğrulama: gerçek H.264 FFmpeg test görüntüsü → simüle eklenti → gerçek Unix soketi → WebRTC/SRTP/RTP → H.264 çözümü otomatik test edildi. Bu, gerçek OBS program videosu, telefon, sahne geçişi veya performans kabulünün yerine geçmez; bunlar hedef bilgisayarda ayrıca doğrulanmalıdır. FFmpeg yalnız isteğe bağlı entegrasyon test aracıdır, uygulama bağımlılığı değildir.

Paylaşılabilir örnek yapılandırma tüm IPv4 arayüzlerini (`0.0.0.0:8080`) kullanır; dosyasız ilk kurulum ise yalnız loopback'te başlar. Yalnızca güvenilir yerel ağda kullanın, internete port yönlendirmesi yapmayın. **HTTP şifreli değildir; ağ trafiğini okuyabilen biri kodu veya oturumu ele geçirebilir.** Güvenilir HTTPS üretim teslim şartı olarak açıktır. Güvenlik duvarı TCP 8080 erişimini engelliyorsa ayrıca izin gerekir; uygulama sistem/güvenlik duvarı ayarlarını değiştirmez.

### Telefon/LAN eşleştirme

1. OBS bilgisayarında `http://127.0.0.1:8080` açıp **Mobil cihaz ile giriş** sekmesine geçin. İsterseniz doğrudan `http://127.0.0.1:8080/pair` de açabilirsiniz. Kendi HTTP portunuz farklıysa onu kullanın.
2. Telefon ve bilgisayar aynı güvenilir ağdayken telefonun normal kamerasıyla QR kodu okutun ve çıkan bağlantıyı açın. Kod elle yazılmadan eşleşip Kontrol ekranına geçilir; tarayıcıya kamera izni vermeniz gerekmez.
3. Birden fazla QR varsa telefonun bağlı olduğu Wi-Fi/LAN adresinin altındakini seçin. LAN adresi kullanılamıyor yazarsa **Ayarlar → HTTP dinleme IP adresi** alanını `0.0.0.0` yapıp kaydedin ve yalnız Controller'ı yeniden başlatın. Güvenlik duvarını tümden kapatmayın.
4. Kamera okuyamıyorsa telefonda `http://SUNUCU_LAN_IP:8080` açıp aynı sayfadaki 12 karakterli kodu elle girin.

QR ve kod 5 dakika geçerli ve tek kullanımlıktır; bir telefonu eşleştirdikten sonra diğer telefon için yeni kod açın. `/pair` sayfasını bilgisayarda yenilemek önceki kodu geçersiz kılar; dördüncü sekme geçerli mevcut kodu paylaşır. Süre dolduğunda veya kod kullanıldığında yeni QR açın. QR yalnız bilgisayardaki loopback erişiminde görünür; eşleşmiş telefon yeni kod/QR göremez veya üretemez.

Oturum 12 saat sürer. **Tüm LAN oturumlarını iptal et** düğmesi veya sunucuyu yeniden başlatmak erişimi kaldırır. Bu nedenle sunucu tekrar açılınca yeniden giriş istenmesi normaldir; tekrar QR okutabilirsiniz. QR kalıcı giriş veya şifreli bağlantı sağlamaz. OBS parolası QR'a konmaz; QR yerelde üretilir, harici bir QR servisine gönderilmez.

Eşleştirilmemiş LAN API istekleri 401 alır; panel sayfaları eşleştirmeye yönlenir. Oturum çerezi HttpOnly/SameSite=Strict, doğrudan TLS altında ayrıca Secure'dür. En fazla 128 oturum ve dakikada 20 başarısız kod denemesi kabul edilir. OBS parolası eşleştirme kodu değildir, tarayıcıya gönderilmez. Host koruması yalnız localhost/loopback ve sunucunun başlangıçta tespit ettiği arayüz IP'lerini kabul eder; keyfi DNS adı ve proxy başlıkları güvenilmez. IP değişince sunucuyu yeniden başlatın. Reverse proxy desteği yoktur; loopback proxy arkasına koymak yerel erişim sınırını geçersiz kılabilir. Tüm istek gövdeleri 5 saniyelik okuma sınırına tabidir.

## Geliştirme

Arayüz 19:9 yatay ekranları önceliklendirir ve mevcut tarayıcı alanına (`100dvh`) sığar; yatay/dikey kaydırma yoktur. Başlık, sekmeler ve alt bilgi yerini korur; içerik kalan alanı kullanır. Dikey ekranlar da desteklenir. Yeni bileşenler eklenirken içerik taşması ayrıca test edilmelidir.

- `main.go`: yalnızca uygulama giriş noktası ve hata çıkışı.
- `Lib/configuration.go`: `Configuration` ve `ReadConfiguration`; JSON okuma ve IP/port doğrulaması.
- `Lib/settings.go`: yalnız loopback'ten bağlantı ayarları, parolanın korunması ve atomik kayıt.
- `Lib/favorites.go`: üç favori sahnenin kalıcı seçimi ve API'leri.
- `Lib/audio.go`: ses kaynakları, doğrulanan seviye ve susturma işlemleri.
- `Lib/obs.go`: kimlik doğrulamalı OBS WebSocket istemcisi, metrikler, geçici screenshot ve sahne değişimi/doğrulaması.
- `Lib/serve.go`: `Run` ve `NewHandler`; Gin yönlendirmesi, şablonlar, statik dosya sunumu ve `http.Server` ile temiz kapanma.
- `web/assets.go`: şablon ve statik dosyaları `go:embed` ile gömer; çalışma dizininden bağımsız sunulur.
- `web/templates/*.tmpl`: ortak düzen ve üç sayfa alanı.
- `web/static/app.css`: sayfa yerleşimi ve Kehribar Gece renkleri.
- `web/static/app.js`: periyodik durum yenileme ve sahne kontrolleri.
- `web/static/audio.js`, `settings.js`: ses işlemleri ve yerel ayar penceresi.
- `web/static/vendor/`: temanın yerel CSS kopyası ve lisansı.

Dosyalar `go:embed` ile derlenir; şablon/CSS değişikliklerinden sonra sunucuyu yeniden başlatın. Node.js, CDN veya çalışma anında internet gerekmez. Henüz PWA yoktur.

```sh
go test ./...
go vet ./...
go build -o obs-lan-controller .
```

Test sonuçları `test_output/` altında tutulur. Gerçek demo OBS üzerinde Deck sahne geçişi, kaynak göster/gizle, sustur/aç ve medya oynat/duraklat/yeniden başlat/durdur döngüsü doğrulandı. Test geçici sahne ve yerel medya girdisi oluşturup kaldırır; önceki program sahnesini geri yükler ve temizliği tekrar okuyarak doğrular. Yayın/kayıt durumu değiştirilmez. Fiziksel telefon/tablet ve uzun süreli yük testi yapılmadı.

2 Kasım 2026 Terra tesliminde Go/Node/gerçek headless Chrome testleri, race kontrolü, vet ve derleme geçti. Yeni derlemeyle gerçek OBS'den durum, ses kaynakları ve favori seçenekleri yalnız okundu. Yeni özelliklerin gerçek sahne/ses komutları, fiziksel cihaz ve kayıt kabulü hâlâ açık: önceki hızlı kayıt testinden kalan kayıt aktif görünüyor. Bu nedenle **Milestone 1'in genel gerçek OBS kabulü verilmedi**; ayrıntılar `PLAN.md` ve test raporlarındadır.

Canlı Deck kabul testi yalnız açık izinle: `OBS_DECK_LIVE=1 go test ./Lib -run '^TestDeckLiveAcceptance$' -v -count=1`. Önce `test_output/deck-acceptance.mp4` yoluna en az birkaç saniyelik sesli yerel video yerleştirin; normal testler OBS'ye bağlanmaz. Görünürlük tuşları koleksiyon ve kaynak kimliğine bağlanır; eski kimliksiz tuşlar düzenleyiciden yeniden atanmalıdır. Düzen yalnız bu tarayıcının localStorage alanındadır.

## Tema kaynağı

[Bootstrap Kehribar Gece](https://github.com/MuhammedOzby/bootstrap-kehribar-gece), commit `0a7d4423c00e3dab021ec70570224f7d41466f4f`, `dist/css/bootstrap.min.css`. Kaynak dosya değiştirilmeden alınmıştır. Lisans: `web/static/vendor/LICENSE.bootstrap`. Uygulamaya özgü yerleşim ve kehribar/gece renkleri `app.css` katmanındadır.
