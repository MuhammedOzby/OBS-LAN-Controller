# Configuration RED

Tarih: 2026-10-29T08:27:38.531549+00:00

Amaç: JSON okuma testi. Ortam: Linux, Go 1.24.1.

Komut: `go test ./Lib -run TestReadConfiguration -v`

```text
# obs-lan-controller/Lib [obs-lan-controller/Lib.test]
Lib/configuration_test.go:14:17: undefined: ReadConfiguration
FAIL	obs-lan-controller/Lib [build failed]
FAIL
```

Çıkış kodu: 1

Gerçek OBS kullanılmadı. Önceki PATH engeli /usr/local/go/bin ile giderildi.
