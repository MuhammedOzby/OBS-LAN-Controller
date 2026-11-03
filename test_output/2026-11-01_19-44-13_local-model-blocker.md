# Yerel model erişimi engeli ve OBS takip durumu

- Tarih: 2026-11-01 19:44:13 +03:00.
- Amaç: alt ajanların yanıt zaman aşımından sonra yerel model altyapısını ve açık test kaydını kontrol etmek.
- Bu bir altyapı/teşhis raporudur; uygulama kabul testi değildir.

## Gerçek sonuçlar

- Ses backend ve Deck tarayıcı teşhis alt ajanları `Timed out waiting for response headers` hatasıyla sonlandı. Başarı kabul edilmedi.
- `ollama ps` (0): `qwen3.8:27b`, 17 GB, GPU, 32768 bağlam, `Stopping...`.
- `free -h` (0): 31 GiB RAM, yaklaşık 18 GiB kullanılabilir; 9.6 GiB swap kullanımda. Bu veri tek başına kök nedeni kanıtlamaz.
- `curl --silent --show-error --max-time 10 http://127.0.0.1:11434/api/version` (0): Ollama 0.40.13; HTTP servisi yanıtlıyor.
- `/api/chat` üzerinden `qwen3.5:9b` için yalnız `Reply with exactly OK.` isteği, `stream=false`, `think=false`, `num_predict=16`, önce `num_ctx=32768`: 120 saniyede hiç yanıt yok, curl 28.
- `ollama stop qwen3.8:27b` (0) ile büyük modelin durdurulması istendi; servis yeniden başlatılmadı.
- Aynı küçük istek `num_ctx=8192`, 90 saniye sınırıyla tekrar denendi: hiç yanıt yok, curl 28. Dolayısıyla yalnız OpenCode istemcisine özgü olduğu söylenemez; doğrudan Ollama üretimi de engelli.
- `curl --silent --show-error --max-time 10 http://127.0.0.1:18080/api/status` (0): gerçek OBS bağlı, yayın kapalı, test kaydı hâlâ `recordActive=true`, `recordPaused=false`, `recordDuration=166`, `recordBytes=542568`. OBS FPS yaklaşık 30. Test kaydının temizliği doğrulanmış değildir.

## Sınırlar ve takip

Kayıt düğmesine kör tekrarlar veya yeni kayıt döngüleri uygulanmadı. OBS/Ollama servisleri zorla öldürülmedi/yeniden başlatılmadı. Kullanıcının profil/kodlayıcı ayarları ve mevcut dosyalar değiştirilmedi. Yerel kod üretimi engellendiği için koordinatör kod yazma yasağını aşmadı veya bulut kod modeline geçmedi. Milestone 1 tamamlanmış değildir. Model görevleri bundan sonra tek tek yürütülmelidir; kayıt durdurma teşhisi önceliklidir.
