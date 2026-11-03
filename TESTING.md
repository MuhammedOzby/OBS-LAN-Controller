# Test ve doğrulama raporlama kuralı

Bu kural OBS LAN Controller projesinin tüm geliştirme süreci boyunca geçerlidir.

- Her test veya test grubu çalıştırmasının sonucunu repo kökündeki `test_output/` klasöründe bir Markdown dosyasına kaydedin.
- Başarısız ve engellenen denemeleri de raporlayın. Önceki raporları silmeyin veya yeni sonuçlarla değiştirmeyin; yeni çalıştırmalar için yeni rapor oluşturun.
- Dosya adı için `YYYY-MM-DD_HH-MM-SS_kisa-aciklama.md` biçimini kullanın; raporda saat dilimini belirtin.
- Her raporda amaç, ortam, çalıştırılan komutlar veya manuel adımlar, gerçek sonuç ve varsa çıkış kodu, sınırlamalar ve takip gerektiren noktalar bulunsun.
- Simülasyon ile gerçek OBS testlerini açıkça ayırın. Gerçek OBS üzerinde denenmeyen davranışları doğrulanmış saymayın.
- Test çalıştırılmadığında bunu açıkça belirtin. Git ve belge kontrollerini uygulama testi gibi sunmayın.
- Parola, token veya başka gizli bilgileri raporlara eklemeyin.

Otomatik testler `go test ./...`, statik analiz `go vet ./...` ile çalıştırılır. JavaScript davranış testleri için Node.js gerekir; uygulamayı çalıştırmak için gerekmez. Gerçek tarayıcı regresyonları için `CHROME_BIN` ile test tarayıcısının yolu verilir. Ortamda tarayıcı yoksa atlanan testleri başarılı tarayıcı doğrulaması olarak raporlamayın. HTTP smoke ve kapanma denemeleri de `test_output/` altında raporlanır.

Normal testlerin çoğu simüle OBS kullanır. Açık izinle çalışan gerçek OBS testleri ayrıca belirtilir; örneğin `OBS_DECK_LIVE=1` geçici kaynaklar oluşturur ve temizler. Gerçek OBS bağlantısı ve kayıt döngüsü için tarihli raporlar vardır; bunlar sonraki kodun veya her kullanıcı bilgisayarının otomatik kabulü değildir. Milestone 1 doğrulama sırası: `docs/milestone-1-acceptance.md`. Başlangıç deposu kontrolleri uygulama testinden ayrı tutulur.

### Gerçek kayıt testi

`OBS_RECORD_LIVE=1 go test ./Lib -run '^TestRecordLiveAcceptance$' -v -count=1` gerçek OBS'de kayıt oluşturur. Yalnız açık kullanıcı izniyle ve kayıt kontrolü test boyunca başka kullanıcı/istemci tarafından değiştirilmeyecekse çalıştırılır. Aktif kayıt varsa test hiçbir komut göndermeden atlanır; atlama başarı kabulü değildir. Başlangıç veya temizlik durumu belirsizse OBS'yi elle kontrol edin. Testin kaydı mutlaka durdurduğunu varsaymayın. Kayıt dosyaları otomatik silinmez.

1–2 Kasım 2026 hızlı döngü denemesinde durdurma takılması gözlendi. Bu kayıt kapanmadan yeni kayıt testi başlatılmaz; son durum `test_output/2026-11-02_00-34-45_obs-readonly-status.md` dosyasındadır. Tekrarlanan eski teşhis testi kaldırıldı; tarihsel raporlar korundu.

Bu oturumdaki taşınabilir test araçları (sistem kurulumu yapılmadı):

```sh
export PATH="/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH"
export CHROME_BIN="/tmp/opencode/chrome-156/chrome-linux64/chrome"
go test ./...
go vet ./...
```

Bu geçici yollar dağıtım gereksinimi değildir; geçici klasör temizlenirse yeniden sağlanmaları gerekir. Normal kullanıcı için kısa OBS ayarları `RECOMMENDATIONS.md` içindedir.
