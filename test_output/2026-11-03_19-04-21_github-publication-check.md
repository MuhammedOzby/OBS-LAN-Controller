# GitHub paylaşım hazırlığı — depo ve belge kontrolü

Tarih: 3 Kasım 2026, 19:04:21 Europe/Istanbul (UTC+03:00).

## Kapsam

Yerel yapılandırma/favoriler ve derlenmiş uygulamanın Git takibinden çıkarılması; .gitignore, paylaşılabilir parola etiketli loopback örneği ve güncelleme yönergesi. İki repoda protokol belgesi eşitlendi. Uygulama mantığı değiştirilmedi.

## Kontroller ve sonuç

- GitHub ana dalı, dosya listeleri ve belgeler bağlantılı GitHub aracıyla okundu. Yerel aktif yapılandırma değiştirilmedi.
- Hazırlanan dosyalar /tmp altında incelendi. Python 3 kontrolleri: örnek JSON ayrıştırma; parola etiketinin bulunması; loopback başlangıcı; iki protokol belgesinin bayt düzeyinde eşitliği; Markdown kod çitleri. Sonuç: geçti, çıkış kodu 0.
- Geçici Git deposunda `git check-ignore --stdin`: configuration.json, favoriler ve Linux/Windows uygulama çıktıları dışlanıyor; configuration.example.json dışlanmıyor. Sonuç: geçti (dışlananlar için 0; örneğin dışlanmadığını gösteren beklenen kod 1).
- `bash -n scripts/install.sh` hazırlanan eklenti betiğinde geçti, çıkış kodu 0. Yalnız kurulum sonu bilgi metni değişti; kurulum/kaldırma çalıştırılmadı.

## Sınırlar

Bu kontroller uygulama testi değildir. Bu teslimde Go/C++ derleme, gerçek OBS, Controller servisi, telefon veya performans testi çalıştırılmadı. Kullanıcının aktif kullanım bildirimi README'de ayrı tutuldu; platform kabulü veya performans garantisi üretilmedi. Commit geçmişi yeniden yazılmaz; ilk yüklemedeki örnek değerlerin eski commitlerde bulunması korunur.
