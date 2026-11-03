# Gerçek OBS: düzeltilmiş kurulumda kayıt döngüsü

- Tarih: 2026-11-01, 18:15–18:17 (+03:00).
- Amaç: kullanıcının düzelttiği OBS/kodlayıcı kurulumunda mevcut Controller kayıt komutlarını yeniden doğrulamak.
- Ortam: aynı bilgisayarda açık gerçek OBS; repoda önceden derlenmiş `./obs-lan-controller` çalıştırıldı. Bu sonuç sonradan eklenecek kodun derleme/regresyon testi değildir.
- Başlatma: `./obs-lan-controller serve -addr 127.0.0.1:18080` (arka plan). Sonunda yalnız bu süreç SIGTERM ile temiz kapatıldı; çıkış başarılı.
- Ön kontrol: `curl --silent --show-error --max-time 8 http://127.0.0.1:18080/api/status` → çıkış 0, `connected=true`, yayın ve kayıt kapalı, `recordPauseSupported=true`, OBS FPS yaklaşık 30.

## Adımlar

Her eylem için şu komut kullanıldı; ACTION sırasıyla aşağıdaki dört değeri aldı:

```sh
curl --silent --show-error --max-time 10 --write-out '\nHTTP %{http_code}\n' \
  --header 'Origin: http://127.0.0.1:18080' \
  --header 'Content-Type: application/json' \
  --data '{"action":"ACTION","confirmed":true}' \
  http://127.0.0.1:18080/api/broadcast
```

| ACTION | HTTP / çıkış kodu | OBS'den tekrar okunan sonuç |
| --- | --- | --- |
| `start-record` | 200 / 0 | `recordActive=true`, `recordPaused=false` |
| `pause-record` | 200 / 0 | `recordActive=true`, `recordPaused=true`, süre 5666 ms |
| `resume-record` | 200 / 0 | `recordActive=true`, `recordPaused=false`, süre 5733 ms |
| `stop-record` | 200 / 0 | `recordActive=false`, `recordPaused=false`, son bildirilen boyut 508238 bayt |

Yayın tüm döngüde kapalı kaldı. Kullanıcının profil/kodlayıcı ayarları değiştirilmedi. Önceki program sahnesi korundu. OBS'nin mevcut kayıt klasöründe kısa test kaydı oluşmuş olabilir; dosya silinmedi, dosya yolu veya oynatılabilirliği bu HTTP API yanıtlarıyla doğrulanmadı.

## Sınırlar

Gerçek OBS testi; simülasyon değildir. Bu bir API kabul denemesidir, tarayıcı düğmeleri/fiziksel telefon doğrulaması değildir. Halka açık yayın başlatma/durdurma denenmedi. Yeni Milestone 1 kodu tamamlanınca ayrıca regresyon gerekir. Komut yanıtı tek başına değil, backend'in OBS tekrar okumasıyla dönen durum kontrol edildi. Uzun kayıt, çıktı dosyasını çözme veya yük testi yapılmadı.
