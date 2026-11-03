# Genel test ve tarayıcı regresyonu

- Tarih: 2026-11-01, yaklaşık 19:12 (+03:00).
- Amaç: taşınabilir Node/Chrome ile mevcut otomatik testleri çalıştırmak.
- Komut: `env PATH="/tmp/opencode/node-v24.22.2-linux-x64/bin:$PATH" CHROME_BIN=/tmp/opencode/chrome-156/chrome-linux64/chrome go test ./... -count=1`
- Çıkış: 1.
- `obs-lan-controller/Lib`: başarılı, 12.020 s. Açık izin değişkeni verilmediğinden gerçek OBS kayıt/Deck testleri bu komutla çalıştırılmadı.
- Ana paket: başarısız; `TestDeckBrowserRegression` 390×844 senaryosunda 3.57 s sonunda boş regresyon sonucu döndürdü (`deck_frontend_test.go:178`). Ana paket toplam 3.954 s. Diğer testlerin tekil sonuçları bu özet çıktıda görünmüyor.
- Web paketi: test dosyası yok.
- Gerçek Chrome süreci kullanıldı; OBS yanıtları tarayıcı fixture/simülasyonudur. Bu bir fiziksel telefon veya gerçek OBS kabulü değildir. Eksik tarayıcı sonucu atlanıp başarı sayılmadı. Sonraki kod değişiklikleri için bu sonuç kabul vermez.
