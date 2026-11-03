# Gerçek OBS: hızlı kayıt döngüsü başarısız

- Rapor zamanı: 2026-11-01 19:10:18 +03:00.
- Amaç: yerel modelin yazdığı `TestRecordLiveAcceptance` ile güncel kaynak kodundan gerçek OBS kayıt döngüsünü doğrulamak.
- Komut: `OBS_RECORD_LIVE=1 go test ./Lib -run '^TestRecordLiveAcceptance$' -v -count=1`
- Sonuç: çıkış 1; test süresi 6.86 s. Başlat/duraklat/devam adımları doğrulandı, hemen ardından durdurma ve temizleme doğrulanamadı.

```text
record active after start-record
record paused after pause-record
record resumed after resume-record
nominal stop-record: could not confirm state (http=502 connected=false); cleanup will retry
final read-back: recording still active (expected inactive after stop)
cleanup: lost OBS connection while confirming stop-record (http=502)
```

## Son durum ve takip

- Gerçek OBS kullanıldı, simülasyon değil. Önceki elle, aralıklı komut testi başarılıydı; bu yeni hızlı döngü başarı sayılmaz.
- Teşhis için mevcut önceden derlenmiş uygulama `./obs-lan-controller serve -addr 127.0.0.1:18080` ile açıldı.
- `curl --silent --show-error --max-time 10 http://127.0.0.1:18080/api/status` (0): OBS bağlı, yayın kapalı, `recordActive=true`, `recordPaused=false`, süre 166 ms, bildirilen kayıt baytı 542568. Bunlar kaydın düzgün ilerlediğini kanıtlamaz.
- Aynı API'ye Origin ve JSON başlıklarıyla `{"action":"stop-record","confirmed":true}` gönderildi: HTTP 502, curl çıkışı 0. HTTP hata kodu başarısızlıktır.
- Kullanıcının önceden aktif kaydı yoktu; açık görünen kayıt bu test tarafından oluşturuldu. Yayın ve profil ayarları değiştirilmedi. OBS kapatılmadı/yeniden başlatılmadı; dosya silinmedi.
- Test kaydının durdurulması ve hızlı geçiş sorunu için yerel modele ayrı, sınırlı teşhis görevi verildi. Bu rapor anında kayıt temizliği henüz doğrulanmadı. Milestone 1 kabulü verilmedi.
