# Backend entegrasyon doğrulaması
Tarih: 2026-10-29T08:51:21.511716+00:00
Ortam Linux. Sahte OBS WebSocket ile auth, screenshot, scene readback, offline, Origin ve bitrate testleri. Gerçek OBS henüz kullanılmadı.
Komut: go test -race ./Lib
```
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1

```
Çıkış: 2
Komut: go vet ./Lib
```

```
Çıkış: 0
Komut: go build -o /root/.hermes/profiles/coder/cache/scratch/obs-controller .
```

```
Çıkış: 0