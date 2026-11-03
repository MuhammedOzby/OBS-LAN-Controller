# OBS LAN Controller — geliştirme planı

Durum: **Milestone 1 tamamlanmadı.** Gerçek OBS metrikleri, sahne kontrolü, LAN eşleştirme, düzenlenebilir Deck ve Controller tarafında WebRTC önizleme bulunur. Yayın/kayıt işlemleri yalnız Yayın Yönetimi'nde doğrudan çalışır. Aşağıdaki eski teslim notları tarihsel kayıt niteliğindedir; güncel açıklar kabul kutularıyla izlenir. İlk platform Linux. Repo adı `obs-lan-controller`.

2 Kasım 2026: kullanıcının isteğiyle Ollama çalışması bırakıldı; GPT-6.1 Terra yarım favori eklerini temizleyip kalıcı favoriler, ses kontrolü ve yerel kurulum akışını testlerle yeniden uyguladı. Yinelenen kayıt testleri birleştirildi; Chrome test hedefi, pencerelerin düzeni ve asenkron arayüz hataları düzeltildi. Tam Go/Node/Chrome testleri, race kontrolü, vet ve derleme geçti (`test_output/2026-11-02_01-04-34_coordinator-final-verification.md`).

Yeni derleme gerçek OBS'den durum, ses kaynakları ve favori seçeneklerini salt okunur başarıyla aldı. **Gerçek sahne/ses komutları ve kayıt kabulü henüz tamamlanmadı:** önceki hızlı testin kaydı hâlâ aktif görünüyor. Yeni kayıt testi başlatılmadı; yalnız koordinatörün eski/yeni geçici Controller sunucuları kapatıldı, OBS'ye dokunulmadı. Ayrıntı: `test_output/2026-11-02_01-04-35_terra-real-obs-readonly.md`. İşaretli kutular uygulama ve belirtilen otomatik doğrulamayı gösterir; genel gerçek OBS kabulü değildir. Basit OBS rehberi `RECOMMENDATIONS.md`, kabul kontrolü `docs/milestone-1-acceptance.md` dosyasındadır.

Son durum (14:43): kullanıcının isteğiyle QR sürümü yeniden başlatıldı ve çalışır bırakıldı. Gerçek OBS okumasında yayın ve kayıt kapalıdır; yukarıdaki aktif kayıt gözlemi 01:04 tarihli rapora aittir. Yeni kayıt döngüsü yapılmadı, hızlı geçiş hatası çözülmüş kabul edilmedi. Rapor: `test_output/2026-11-02_14-43-12_qr-restart-smoke.md`.

Test raporlama: Proje boyunca her test veya test grubu çalıştırmasının sonucu `test_output/` altında Markdown raporu olarak saklanır. Başarısız/engellenen denemeler de kaydedilir. Ayrıntılar: `TESTING.md`.

## Yeni istek — QR ile mobil giriş

Kullanıcının 2 Kasım 2026 isteği, önceki tam üç sekme kuralını değiştirir: dördüncü sekme **Mobil cihaz ile giriş** olacaktır. Bağlantı ayarları ayrı sekmeye taşınmaz; ortak yerel ayarlar penceresinde kalır.

- [x] `/pair` üzerinde yerelde üretilen QR ve elle kod girişi yedeği.
- [x] Dördüncü sekmede aynı güvenli eşleştirme akışı; yalnız gerçek loopback istemcisine kod/QR gösterme.
- [x] Dinleyicinin gerçek LAN adresi/portuyla, 5 dakikalık tek kullanımlık QR bağlantısından otomatik giriş. Token sorguda değil URL fragment'inde taşınır; tarayıcı geçmişinden temizlenip aynı Origin POST ile eşleştirilir.
- [x] Kodun süre sonu/tekrar kullanım/iptal ve Origin/Host korumaları; bağımsız QR çözme, gerçek loopback/LAN soket akışı ve headless Chrome testleri. Dar/yatay ekran yerleşimleri doğrulandı.
- [ ] Fiziksel telefon kamerasıyla okutma ve kullanıcının gerçek ağında kabul.

QR bu istekle öne alınmıştır; güvenilir HTTPS, kalıcı oturum veya tam PWA teslimi anlamına gelmez. Oturumlar yine 12 saatlik ve süreç içidir; yeniden başlatmada silinir. Test/race/vet/derleme geçti; QR için gerçek OBS kullanılmadı, kullanıcı sunucusu yeniden başlatılmadı. Son fiziksel cihaz kabulü açık tutulur. Rapor: `test_output/2026-11-02_14-40-11_qr-coordinator-verification.md`.

## Ürün ve mimari

Zorunlu arayüz teması: **Bootstrap Kehribar Gece** — https://github.com/MuhammedOzby/bootstrap-kehribar-gece. `.tmpl` tabanlı panel bu temayı kullanacaktır.

Go CLI OBS bilgisayarında çalışır ve `.tmpl` tabanlı paneli sunar. Telefon aynı LAN'dan bu sunucuya bağlanır. Kontrol/metrikler Go üzerinden OBS WebSocket'e; video ise ayrı OBS eklentisinden Go'ya, oradan WebRTC ile tarayıcıya gider.

Paketleme hedefi kullanıcıya tek kurulum deneyimidir; Go uygulaması ve yerel OBS modülü ayrı derleme çıktılarıdır. Windows/macOS desteği ayrı paket ve test gerektirir. Güncel OBS Studio kullanılır; ayrı obs-websocket kurulmaz.

## Milestone 1 — Hermes'in ilk görevi

- [x] Go proje iskeleti, `serve` komutu, temiz kapanma ve gömülü `.tmpl`/CSS/JS dosyaları.
  - Kullanıcı tercihiyle yönlendirme ve statik dosya sunumu Gin Web Framework'e taşındı; mevcut şablonlar ve temiz kapanma korundu.
  - Sunucu ve handler `Lib/serve.go` içine ayrıldı; `main.go` yalnızca giriş noktasıdır. Gömülü varlıklar `web/assets.go` içinde tutulur.
  - Tasarım ekranlarının LAN incelemesi için kullanıcı isteğiyle varsayılan dinleme adresi `0.0.0.0:8080` yapıldı. Bu, eşleştirme/oturum güvenliğinin tamamlandığı anlamına gelmez; yalnızca güvenilir LAN'da kullanılmalıdır.
- [x] Bootstrap Kehribar Gece temasıyla uyarlanabilir arayüz; güncel kullanıcı tercihiyle dört sekme.
  - Kontrol: program görüntüsü/bağımsız video durumu, metrikler, üç favori yuva ve ses paneli. Ortak yerel ayarlar ve favori pencereleri etiketli, temalı ve kısa ekranlarda kaydırılabilir. Headless Chrome yatay/dikey yerleşim ve etkileşim testleri geçti; fiziksel telefon/tablet doğrulanmadı.
- [x] Sunucu tarafında OBS adresi/port/parola yapılandırması ve bağlantı kontrolü.
  - `configuration.json` ile OBS IP/port/parola ve HTTP dinleme IP/port ayarları okunur; okuma ve doğrulama `Lib/configuration.go` içindedir. `Lib/obs.go` kimlik doğrulamalı bağlantıyı kurar; tarayıcıya yalnız güvenli durum projeksiyonu iletilir. Gerçek OBS bağlantısı doğrulandı. Yerel parola dosyası Git dışında tutulur.
- [x] İlk kurulum için yerel erişim sınırı; LAN paneli için temel eşleştirme ve oturum koruması.
  - Dosyasız başlangıç yalnız loopback'te; yerel ortak ayarlar penceresi parolayı geri döndürmez. Atomik `0600` kayıt, başarısız yazmada canlı ayarların korunması, parolayı koruma/değiştirme/silme ve bootstrap dosyasını ezmeme test edildi.
  - LAN eşleştirme: 5 dakikalık tek kullanımlık kod, 12 saatlik HttpOnly/SameSite oturum, iptal ve deneme sınırı. Gerçek socket/Host/Origin denetimleri ve eşleşmiş LAN'dan yönetim erişiminin reddi üretim middleware zinciriyle test edildi. Önceki gerçek ağ eşleştirme raporları korunur. HTTP şifreli değildir; reverse proxy desteklenmez. Güvenilir HTTPS/PWA ayrı aşamadır.
- [x] OBS bağlantı/kimlik doğrulama hataları için güvenli sabit mesajlar ve kontrollü yeniden bağlantı.
  - Hata yanıtlarında ham bağlantı hatası veya parola yoktur; kopma sonrası eski metrikler temizlenir, kritik komutlar yeniden oynatılmaz. Otomatik regresyonlar geçti.
- [x] Sahne listesinden tam üç favori yuva seçimi, ayarları kaydetme ve aktif sahne gösterimi.
  - İlk kayıttan önce ilk üç sahne kullanılır; kaydedilen boş/silinmiş yuva sessizce yeniden atanmaz. Koleksiyon değişimi doğrulanır; seçimler yapılandırmanın yanında atomik saklanır. Komut yetkisi taze OBS durumuna dayanır. Disk hatası/kalıcılık, yetki, mock sahne değişimi, LAN ve Chrome düzenleyici testleri geçti. Yeni derlemeyle gerçek OBS yalnız okundu; gerçek sahne komutu kabulü açık.
- [x] Ses kaynaklarında seviye ayarlama ve susturma; harici OBS değişikliklerinin panele yansıması.
  - Ses destekli kaynak keşfi, %0–100 seviye, açık susturma değeri, sınırlanmış/sıralı kaydırıcı istekleri ve OBS tekrar okuması eklendi. Harici güncelleme, hatalı yanıt, bağlantı kopması ve son kaydırıcı değerinin korunması mock/Node ile; arayüz Chrome ile doğrulandı. Gerçek kaynak listesi okundu; gerçek seviye/susturma komutu kabulü açık.
- [x] Saniyede yaklaşık bir metrik örneği: OBS CPU/bellek, aktif FPS, render süresi, render/kodlama atlamaları ve yayın çıkışı durumu.
- [x] Gönderilen baytların zaman farkından bit hızı hesabı; sayaç sıfırlanması ve yayın kapalı durumu ele alınır.
  - İlk örnekte değer bilinmiyor olarak kalır. Gerçek OBS'de ardışık örneklerle, pasif yayın/sayaç sıfırlanması birim testleriyle doğrulandı. Uzun süreli yük testi yapılmadı.
- [ ] Kritik yayın/kayıt işlemleri yalnızca üçüncü sekmede; doğrudan tek tıklama, işlem sırasında kilitleme ve gerçek durum takibi.
  - Uygulama ve mock/Node korumaları mevcut: açık başlat/durdur/duraklat/devam komutları, tek tıklama, işlem kilidi, Origin kontrolü ve taze durum doğrulaması. Ek onay penceresi yoktur. Duraklatma desteği OBS profiline bağlıdır.
  - Düzeltilmiş OBS kurulumunda aralıklı gerçek kayıt döngüsü geçti; hızlı döngüde durdurma takıldı. Son gerçek okumada test kaydı aktif. Bu sorun çözülmediğinden birleşik kabul kutusu açık tutulur; yeni kayıt/yayın testi yapılmadı.
- [x] Deck ve önizlemenin gerçek bağlantı/kapsam durumunu açık gösterme.
  - Deck artık iskelet değil, ikinci aşamadaki uygulamadır. Önizleme WebRTC kullanır; ekran görüntüsü döngüsü nihai video diye sunulmaz. Video ve OBS kontrol bağlantısı bağımsızdır. Gerçek video/eklenti kabulü üçüncü aşamada kalır.
- [x] Kurulum/çalıştırma ve doğrulama adımlarını anlatan kısa README.

Kabul: Gerçek OBS bağlantısında sahne/ses kontrolü ve metrikler çalışır; bağlantı kopması yanlış durum veya tekrar gönderilen kritik işlem üretmez. İlk iki sekmede başlat/durdur yoktur. Önizleme veya Deck tamamlanmış gibi sunulmaz. Gerçek OBS kullanılamadıysa bu kabul henüz sağlanmamış olarak raporlanır. Bu aşamada eklenti, video aktarımı ve tam PWA kurulumu teslim şartı değildir.

## Milestone 2 — özelleştirilebilir Deck

- [x] Tuş ekleme/düzenleme/silme, sıralama, etiket, ikon ve renk.
- [x] Sahne değiştirme, kaynak gösterme/gizleme, susturma ve medya kontrolü.
- [x] Gerçek OBS durumuna bağlı tuş görünümü ve kalıcı düzen.
- [x] İzinli eylem listesi; kritik işlemleri veya genel komut/makroyla aynı etkiyi eklemeye izin verilmez.
  - 12 yuvalı Deck düzeni bu tarayıcıda localStorage ile kalıcıdır; cihazlar arası eşitleme yoktur. Kaybolmuş hedef sessizce değiştirilmez; görünürlük tuşları koleksiyon/sahne/kaynak kimliğine bağlanır, kimliksiz eski tuşlar yeniden atanmalıdır. Portre grid çakışması düzeltildi; 390×844, 844×390, 1280×720 gerçek Chromium regresyonları geçti. Gerçek demo OBS'de geçici sahne ve sesli medya ile sahne, görünürlük, susturma ve tüm medya eylemleri çalıştırıldı; restart medya imleciyle doğrulanır. Geçici kaynakların yokluğu ve eski program sahnesine dönüş tekrar okunarak doğrulandı; yayın/kayıt değiştirilmedi. Raporlar: `test_output/*_deck-cursor-fixed-live.md`, `*_deck-combined-regressions.md`. Fiziksel cihaz ve uzun süreli yük testi yapılmadı. LAN güvenliği ayrı ilk teslim/dağıtım koşuludur.

## Milestone 3 — canlı önizleme entegrasyonu

İlk Controller entegrasyonu uygulandı: Linux aynı-UID Unix soketi, OLCP0001 kayıt doğrulama, Pion H.264 WebRTC (yeniden kodlama yok), sınırlı izleyici/GOP kuyrukları, heartbeat/S/T/K ve görünür sayfa yaşam döngüsü. Panel artık screenshot istemez; video ve OBS kontrol bağlantısı bağımsızdır. Otomatik testte FFmpeg baseline H.264 test görüntüsü simüle eklentiden gerçek Unix soketi/WebRTC/SRTP üzerinden alınıp yeniden çözüldü; gerçek OBS programı, fiziksel telefon ve performans ölçümü **henüz yapılmadı**. Bu nedenle birleşik kabul kutuları açık bırakıldı. Kaynak FPS gösterimi encoder yapılandırmasıdır, çözülmüş FPS/gecikme ölçümü değildir. Aynı bilgisayarda çalıştırma ve Arch bağımlılıkları README'de. Başlangıç checkpoint: `077edf1`; wire formatı değiştirilmedi, eklenti tekrar kurulumu gerekmiyor.

- [ ] `docs/preview-protocol.md`: eklentiyle sürümlü, doğrulanmış yerel video bağlantısı sözleşmesi.
- [ ] Önizleme prototipi: codec/encoder seçimi, platformda bulunabilirlik ve tarayıcı uyumluluğu.
- [ ] Eklentiden kodlanmış videoyu alıp yeniden kodlamadan WebRTC ile taşıma.
- [ ] Yayına giden program çıktısı; başlangıç hedefi 720p/30 FPS. Bu bir ölçüm hedefidir, garanti değildir.
- [ ] İzleyici yokken önizleme kodlamasını/aktarımı durdurma; tekrar açılışta anahtar kareyle başlama.
- [ ] OBS bağlantısı ve video bağlantısı ayrı durum göstergelerine sahip olur.
- [ ] Önizleme FPS ve mümkün olan aktarım istatistikleri; gecikme değeri ancak güvenilir ölçülebiliyorsa gösterilir.
- [ ] Önizleme kopması kontrol panelini bozmaz; Go servisi kapanırsa OBS eklentisi OBS'yi engellemez.

## Milestone 4 — PWA ve Linux dağıtımı

- [ ] Manifest, ikonlar, kapsamı sınırlı service worker ve çevrimdışı uygulama kabuğu.
- [ ] Telefonlarda güvenilir HTTPS ve sertifika güveni için uygulanabilir kurulum akışı; tarayıcı kontrollerini kapatmayı gerektirmez.
- [x] Eşleştirme QR kodu, oturum iptali ve gizli bilgilerin önbelleğe alınmaması (2 Kasım isteğiyle öne alındı; fiziksel telefon kabulü yukarıda açık).
- [ ] `setup` ve `doctor` akışlarının tamamlanması; hata teşhisi ve desteklenen OBS sürüm/özellik denetimi.
- [ ] İlk desteklenen Linux dağıtımı, mimari ve OBS kurulum biçiminin açıkça seçilmesi; ardından diğer kurulum biçimleri.
- [ ] Flatpak kullanılıyorsa eklenti dağıtımı, süreçler arası bağlantı ve ağ izinlerinin ayrıca doğrulanması.
- [ ] İsteğe bağlı kullanıcı servisi ve otomatik başlatma; kurulum yetkisi dışında sistem değişikliği yapılmaz.

## Performans kabulü

Aynı OBS sahnesi/çıkış ayarlarında üç koşulu karşılaştırın: eklenti kapalı, eklenti açık ama önizleme izleyicisi yok, tek telefon izleyicisi var. OBS CPU, render süresi, render/kodlama atlamaları; Controller CPU/bellek ve önizleme akıcılığı ölçülür. Mümkünse GPU/encoder yükü ayrıca sistem araçlarıyla ölçülür; standart OBS metriği gibi sunulmaz.

Nicel ek yük bütçesi ilk Linux donanım ölçümünden sonra belirlenir. Sürekli büyüyen kuyruk/bellek, OBS'yi bekleten aktarım veya belirgin yayın bozulması kabul edilmez. Sorunda önizleme kalitesi düşürülebilmeli veya önizleme kapatılabilmelidir.

## Kapsam dışı / sonraki işler

İlk sürümde internetten uzaktan erişim, izleyici analitiği, OBS screen toolset, bulut hesabı ve platform sohbet entegrasyonları yoktur. Windows/macOS ortak kodla hedeflenir fakat Linux testleri bu platformlarda destek kanıtı sayılmaz.


## 2 Kasım 2026 15:05 — manuel liste denemesi

Gerçek OBS’ye bağlı mevcut Chrome panelinde boş parola ile ayar kaydı/yenileme, üç favori/boş yuva kalıcılığı, gerçek sahne komutu ve başlangıca dönüş, Desktop Audio seviye/susturma komutları ve geri yükleme doğrulandı. LAN adresindeki ayrı Chrome sekmesinde eşleşmemiş erişim engeli, kodla giriş ve yönetim kontrollerinin gizlenmesi gözlendi. OBS yerel pencere erişimi olmadığından OBS içinden ters yön ses, WebSocket sunucu kesintisi ve kayıt/dosya oynatma denenmedi. Controller yeniden başlatılmadı; çalışan kaynak sürümü doğrulanmadı. Fiziksel telefon/QR ve mobil yerleşim kabulü açık; ekran boyutu denemesinde tarayıcı aracı koptu, oturum iptali ve ekran boyutunu geri alma doğrulanamadı. Ara ses geri yüklemesinde ilk susturma tıklaması uygulanmamış göründü, ayrı tekrar ile başlangıç durumu doğrulandı; nedeni açık. Doğrudan /pair gezinmesinde Forbidden Origin görüldü, panel içi bağlantı çalıştı; kesin kök neden ölçülmedi. Tam kabul kutuları değiştirilmedi. Ayrıntı: `test_output/2026-11-02_15-05-00_manual-checklist-browser.md`.

## 3 Kasım 2026 — GitHub paylaşım hazırlığı

- [x] Yerel yapılandırma, favoriler ve derlenmiş uygulama kaynak takibinden çıkarıldı; .gitignore eklendi.
- [x] Örnek yapılandırmaya parola etiketi ve loopback başlangıcı eklendi; indirme ve ilk yüklemeden güncelleme adımları açıklandı.
- [x] Kullanıcının aktif kullanım bildirimi README'ye, tarihli otomatik test raporlarından ayrı kaydedildi.

Bu teslim uygulama mantığını değiştirmez. Depo/belge kontrolleri uygulama veya gerçek OBS testi değildir; performans ve platform kabul kutuları bu işlemle kapanmaz. İlk commit geçmişi korunur; başlangıçtaki örnek değerler geçmişten silinmez.

## 3 Kasım 2026 — İngilizce ve Türkçe README

- [x] Ana README.md İngilizceye çevrildi; mevcut Türkçe içerik README-TR.md olarak korundu.
- [x] Her iki dilin başına AI desteğiyle geliştirme notu ve karşılıklı dil bağlantısı eklendi.
- [x] Komut blokları, Markdown kod çitleri ve göreli belge bağlantıları kontrol edildi.

Yalnız belgeler değiştirildi; uygulama, kurulum ve gerçek OBS testleri yeniden çalıştırılmadı.
