# Bütünleşik dashboard doğrulaması
Tarih: 2026-10-29T08:55:58.790193+00:00
Linux Go ve Node; mock OBS ve frontend DOM harness. Bu grup gerçek OBS testi değildir. Race daha önce CGO/derleyici eksikliği nedeniyle engellendi.
Komut: go test ./... -count=1
```
ok  	obs-lan-controller	0.125s
ok  	obs-lan-controller/Lib	0.009s
?   	obs-lan-controller/web	[no test files]

```
Çıkış: 0
Komut: go vet ./...
```

```
Çıkış: 0
Komut: node --check web/static/app.js
```

```
Çıkış: 0
Komut: go build -o /root/.hermes/profiles/coder/cache/scratch/obs-controller .
```

```
Çıkış: 0