# frontend-full-suite

UTC: 2026-10-29T08:54:46.856944+00:00

Komut: /usr/local/go/bin/go test ./... && /usr/local/go/bin/go vet ./... && node --check web/static/app.js

Çıkış: 0

```text
ok  	obs-lan-controller	0.177s
ok  	obs-lan-controller/Lib	0.010s
?   	obs-lan-controller/web	[no test files]
```

Sınır: Gerçek OBS bağlantısı veya mutasyon yapılmadı. BrowserServer normal çalışmada bilinçli skip; gerçek Chromium testleri ayrı raporda. Node frontend davranış testleri için geliştirme ortamında gerekir; uygulama çalışma zamanı bağımlılığı değildir.
