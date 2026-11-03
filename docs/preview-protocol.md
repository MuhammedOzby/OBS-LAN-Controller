# Yerel önizleme protokolü — taslak v1

Durum: Eklenti ve Controller deneysel v1 wire formatını uygular. Controller Linux Unix soketini sunar ve H.264 verisini yeniden kodlamadan WebRTC ile tarayıcıya taşır. 3 Kasım 2026'da proje sahibi aktif kullanım bildirmiştir; sürümlü OBS/cihaz kabulü ve sayısal performans ölçümleri ayrıca kaydedilmelidir. Otomatik taşıma testleri simüle gönderici ve gerçek WebRTC kullanır.

## Güven sınırı ve taşıma

Linux AF_UNIX/SOCK_STREAM, `${XDG_RUNTIME_DIR}/obs-lan-controller/preview.sock`. Controller sunucu, eklenti istemcidir. Controller dizini 0700 ve soketi 0600 oluşturmalı; eski soketi yalnız kendi sahipliğini doğruladıktan sonra temizlemelidir. İki taraf `SO_PEERCRED` ile aynı etkin UID'yi doğrular. Eklenti dizin/soket için `lstat` sahiplik/tür/izin kontrolü yapar, symlink ve diğer kullanıcılara açık hedefleri reddeder. Controller da peer UID'sini doğrulamadan oturum kabul etmemelidir. Eklentinin yerel taşıması TCP/LAN video portu açmaz ve OBS parolası taşımaz. Controller'ın tarayıcıya sunduğu HTTP/WebRTC bağlantısı bu yerel soketten ayrıdır.

Bu prototipte yetki sınırı aynı Unix kullanıcısıdır; aynı UID altında kötü niyetli süreçlere karşı uygulama kimliği/ayrı oturum sırrı sağlamaz. Ayrı kullanıcıdaki Controller desteklenmez. Daha güçlü tehdit modelinde nonce/HMAC veya doğrulanmış fd aktarımı ayrı sürümde tasarlanmalıdır. Flatpak dosya sistemi/UID/runtime görünürlüğü ayrıca doğrulanmalıdır; `--filesystem=host` önerilmez.

## El sıkışma ve kontrol

1. İstemci sekiz ASCII baytı `OLCP0001` yollar.
2. Sunucu aynı sekiz baytla yanıtlar. En fazla 500 ms; farklı sürüm/kısmi/yanıtsız bağlantı kapanır.
3. Sunucu tek baytlık komutlar yollar:
   - `S`: ilk izleyici; bağımsız önizleme encoder/output oluştur ve başlat. Tekrarlanması yeni nesil başlatır.
   - `T`: son izleyici ayrıldı; output/encoder bırak, bağlantı boşta kalabilir.
   - `K`: anahtar kare isteği; mevcut ilk sürüm encoder'ı yeniden oluşturur (geçici boşluk olabilir).
   - `P`: heartbeat, önerilen aralık 1 saniye.
4. Beş saniye kontrol baytı gelmezse kaynaklar bırakılır ve bağlantı kapanır. Bilinmeyen komut bağlantıyı kapatır. Yeniden bağlanma en az iki saniye aralıklıdır. Başlangıçta/yeniden bağlanmada kendiliğinden encoder açılmaz: yeni `S` gerekir.

Komutlar güvenilir yerel kanalda yalnız Controller'dan gelir. Controller ilk/son yetkili izleyici sayısını belirlemeli; tarayıcı doğrudan sokete bağlanmaz. Controller her bağlantıda tek gönderici kabul etmeli ve yavaş tüketicide okumayı sonsuza kadar bekletmemelidir.

## Eklenti → Controller kayıtları

Tüm sayılar big-endian. TCP varsayımı yoktur; stream okumaları kısmi olabilir. Önce 4 bayt uzunluk `N` (bu alan hariç), sonra tam N bayt gövde. Üst sınır: `N <= 4*1024*1024 + 36`. Kontrol/başlıklar ayrıştırılmadan bellek ayrılmaz.

| Gövde ofseti | Boyut | Anlam |
| --- | --- | --- |
| 0 | 1 | Tür: 2=konfigürasyon, 3=video access unit |
| 1 | 1 | bit0=anahtar kare, diğer bitler 0 |
| 2 | 2 | ayrılmış, 0 |
| 4 | 4 | nesil: her başarılı S/K sonrası artar, bağlantı içinde 1'den başlar |
| 8 | 8 | signed int64 PTS (ikiye tümleyen) |
| 16 | 8 | signed int64 DTS (ikiye tümleyen) |
| 24 | 4 | timebase payı (pozitif) |
| 28 | 4 | timebase paydası (pozitif) |
| 32 | 4 | payload uzunluğu, `N-36` ile eşit |
| 36 | değişken | payload |

Tür 2: payload `width:u32, height:u32, fps_num:u32, fps_den:u32` ve ardından OBS x264'ün verdiği Annex-B H.264 SPS/PPS ekstra verisi (en fazla 64 KiB). Zaman alanları 0, timebase 1/30, bayrak 0. Konfigürasyon gerçek encoder başlatıldıktan sonra ve o neslin ilk video kaydından önce gönderilir. Tür 3: OBS x264 Annex-B access unit; zamanlar encoder'ın kendi timebase'indedir. PTS/DTS'yi yerel duvar saati veya uçtan uca gecikme olarak yorumlamayın. Ses yoktur. H.264 baseline, B-frame'siz düşük gecikme ayarları; bu sürümde codec pazarlığı/hardware seçimi yoktur.

## Sınırlar ve toparlanma

Kuyruk en fazla 60 paket / 4 MiB; tek access unit en fazla 4 MiB. Encoder callback'inde ağ/disk yoktur. Kapasite aşılırsa tüm bekleyen GOP atılır ve sonraki doğal anahtar kareye kadar delta kare kabul edilmez; bağımlı kareler rastgele bırakılarak bozuk akış sürdürülmez. Keyframe aralığı hedefi bir saniyedir. Gönderim 500 ms içinde tamamlanmazsa bağlantı kapanır ve encoder/output bırakılır. Başlangıçta, taşmadan sonra ve yeniden bağlantıda ilk aktarılan video anahtar karedir. Controller yeni nesilde eski decoder/RTP durumunu sıfırlamalı, SPS/PPS'yi anahtar kareyle birlikte alıcıya sağlamalıdır. Kayıtlar akıştan düşürülmez; yavaş peer için bağlantı sıfırlanır.

## Henüz uygulanmayan sözleşme parçaları

Yapılandırılmış hata/istatistik mesajı, codec pazarlığı, H.264 dışındaki codec'ler, gerçek donanım encoder seçimi, uygulama düzeyinde ACK bu yerel protokolde yoktur. İzleyici sınırı ve tarayıcı WebRTC sinyalleşmesi Controller tarafından uygulanır. Eklenti başlatma hatasını OBS loguna yazar ve soketi kapatır; Controller bunu önizleme bağlantı hatası olarak göstermeli, OBS kontrol bağlantısını kesmemelidir.
