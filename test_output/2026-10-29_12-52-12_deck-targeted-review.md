# Deck odaklı test

Tarih/saat (UTC): 2026-10-29T12:52:12.280160+00:00

Amaç: Deck backend fixture ve frontend model/şablon testleri. Gerçek OBS mutasyonu yok.

Komut: `/usr/local/go/bin/go test ./... -run TestDeck -count=1`

Çıkış: 0
```text
ok  	obs-lan-controller	0.050s
ok  	obs-lan-controller/Lib	4.116s
?   	obs-lan-controller/web	[no test files]

```

Sınırlamalar: gerçek OBS ve fiziksel mobil cihaz üzerinde doğrulama yapılmadı; DOM davranışlarının tümü bu testte kapsanmıyor.
