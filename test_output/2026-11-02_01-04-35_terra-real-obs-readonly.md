# Yeni Terra derlemesiyle gerçek OBS salt okunur kontrolü

- Zaman: 2026-11-02 01:04 +03:00.
- Amaç: yeni derlemenin gerçek OBS durum, ses kaynağı ve favori seçeneklerini okuyabildiğini doğrulamak; hiçbir sahne/ses/kayıt/yayın komutu göndermemek.
- Başlatma: `/tmp/opencode/obs-lan-controller-terra-final serve -addr 127.0.0.1:18081`.
- Yapılandırma mevcut yerel dosyadan sunucu tarafında okundu; parola çıktılanmadı, dosya değiştirilmedi.
- Eski test sunucusu önizleme soketini tuttuğundan yeni sunucu `preview already running` bildirdi; HTTP/OBS kontrolü devam etti. Bu çalıştırma video kabulü değildir.

Aşağıdaki her istekte `curl --silent --show-error --max-time 10 --write-out '\nHTTP %{http_code}\n' URL` kullanıldı; hepsi curl çıkış 0 ve HTTP 200:

| URL | Gerçek sonuç |
| --- | --- |
| `http://127.0.0.1:18081/api/status` | OBS bağlı; yayın kapalı; önceki test kaydı aktif, duraklatılmamış; süre 166 ms ve boyut 542568 bayt; yaklaşık 30 FPS; varsayılan üç favori döndü |
| `http://127.0.0.1:18081/api/audio` | Ses destekli dört gerçek kaynak ve hacim/susturma durumları döndü; biri susturulmuştu |
| `http://127.0.0.1:18081/api/favorites` | Güncel koleksiyon ve üç mevcut sahne/seçenek döndü; kaydetme yapılmadı |

Yalnız bu oturumun yeni 18081 süreci SIGTERM ile kapatıldı (çıkış 0). Ardından önceki oturumda koordinatörün başlattığı eski 18080 test sunucusu da tam komut eşleşmesiyle SIGTERM kullanılarak temiz kapatıldı (çıkış 0); eski derleme/önizleme soketi çalışır bırakılmadı. OBS ve Ollama süreçlerine dokunulmadı. Controller bağlantısını kapatmak OBS kaydını durdurmaz; test kaydının temizlendiği iddia edilmez.

Gerçek OBS, simülasyon değil. Sahne değiştirme, favori kaydetme, ses değiştirme, kurulum ayarı yazma, kayıt/yayın, fiziksel cihaz veya gerçek video testi yapılmadı. OBS'nin aktif görünen test kaydı hâlâ açık bir konudur; Milestone 1 genel gerçek OBS kabulü verilmedi.
