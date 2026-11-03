# Yayın yönetimi entegre test
UTC: 2026-10-29T09:35:17.973925+00:00
Ortam Linux; Go/Node ve mock OBS. Gerçek çıkış değiştirilmedi.

Komut: `go test ./... -count=1`
Çıkış: 0
```text
ok  	obs-lan-controller	0.195s
ok  	obs-lan-controller/Lib	0.550s
?   	obs-lan-controller/web	[no test files]

```

Komut: `go test -race ./...`
Çıkış: 2
```text
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1

```

Komut: `go vet ./...`
Çıkış: 0
```text

```
