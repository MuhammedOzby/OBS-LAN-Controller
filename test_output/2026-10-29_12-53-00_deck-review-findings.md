# Deck son inceleme bulguları

Tarih: 2026-10-29, UTC. Bu belge kod incelemesi özetidir; ek bir test çalıştırması değildir. Ayrıntılı gerçek çalıştırmalar aynı klasördeki deck-final-go-review, deck-targeted-review ve deck-browser-review raporlarındadır.

## Tarayıcı fixture ile yeniden üretilenler

- Portre 390x844: grid sabit 4x108px satırlara rağmen ebeveyn tarafından küçülüyor; tuşlar alt içerikle çakışıyor. web/static/deck.css:3,93.
- Silinmiş sahne hedefini düzenlemek ilk mevcut hedefi sessizce seçiyor; yalnız etiket kaydetmek hedefi değiştiriyor. web/static/deck.js:139-145.

## Kaynak kodla doğrulanan, uçtan uca istismar/OBS testi yapılmayanlar

- Yüksek: Origin/Host eşitliği güvenilir Host veya oturum doğrulaması değildir. Lib/deck.go:401-408, Lib/serve.go:39-56. DNS rebinding tehdidi mevcut; LAN erişimi eşleştirme/oturum olmadan güvenli kabul edilmemeli. Bu genel sunucu güvenliği açığıdır, yalnız Deck özelliğine özgü değildir.
- Orta: media-restart doğrulaması sadece PLAYING durumunu kontrol eder; zaten oynayan medyanın gerçekten başa sarıldığı kanıtlanmaz. Lib/deck.go:178-179,277-279,317-322.
- Orta: görünürlük hedefi sahne adı+sceneItemId ile eşleşir; koleksiyon değişimi/kimlik yeniden kullanımı farklı kaynağa komut yöneltebilir. Koleksiyon kimliği veya kaynak kimliği ile bağlama yok. Lib/deck.go:150-156,193-196,242.
- Orta: JSON gövdesi çözümlemesi 8s context öncesindedir ve sunucuda body read deadline/ReadTimeout yoktur. Boyut sınırı yavaş gövde aktarımını süre bakımından sınırlamaz. Lib/deck.go:410-420, Lib/serve.go:149-153.

## Test ve teslim durumu

- Deck odaklı Go/Node testleri geçti; go vet geçti.
- Genel go test ./... başarısız. Deck editor dialog/header için aşırı geniş eski test yasakları ve eksik selector için varmış gibi element döndüren DOM harness yeni Deck guard'ını tetikliyor (frontend_test.go:144; app.js:5). Üretim davranışı bozuk olduğu sonucuna yalnız bu hatalarla varılamaz; testler düzeltilmeli.
- Race testi GCC eksikliği nedeniyle engellendi.
- Eşleştirme/oturum, gerçek OBS Deck işlemleri ve fiziksel telefon doğrulaması eksik. PLAN Milestone 2 kutuları açık.
- Uygulama kodu değiştirilmedi; gerçek OBS mutasyonu yapılmadı. İlk iki bulgu DOM/fixture ile yeniden üretildi, diğerleri statik kod incelemesidir.
