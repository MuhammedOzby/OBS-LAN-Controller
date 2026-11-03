# Backend RED
Tarih: 2026-10-29T08:50:23.123564+00:00
Komut: go test ./Lib
Sentetik test; gerçek OBS kullanılmadı.
Beklenen: yeni API sembolleri henüz yok.
```
FAIL	obs-lan-controller/Lib [build failed]
FAIL
# obs-lan-controller/Lib [obs-lan-controller/Lib.test]
Lib/obs_test.go:13:12: undefined: NewOBSClient
Lib/obs_test.go:14:18: undefined: newHandler
Lib/obs_test.go:32:9: undefined: bitrate
Lib/obs_test.go:39:12: undefined: NewOBSClient

```
Çıkış kodu: 1