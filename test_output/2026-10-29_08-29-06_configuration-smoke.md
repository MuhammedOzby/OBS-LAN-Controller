# Yapılandırma HTTP smoke ve race


Tarih: 2026-10-29T08:29:06.586118+00:00

Ortam: Linux, Go 1.24.1. Gerçek OBS kullanılmadı.


## `CGO_ENABLED=1 go test -race ./...`
```text
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
FAIL	obs-lan-controller [build failed]
FAIL	obs-lan-controller/Lib [build failed]
?   	obs-lan-controller/web	[no test files]
FAIL
```
Çıkış kodu: 1


HTTP smoke: farklı çalışma dizininden -config ile başlatılan derlenmiş uygulama, dosyadaki loopback IP/port üzerinde GET / için HTTP 200 ve Kontrol sayfası döndürdü. Yapılandırma dosyası URLleri HTTP 404 döndürdü.

SIGTERM temiz kapanma çıkış kodu: 0

Sınırlama: OBS bağlantısı uygulanmadı; fiziksel cihaz doğrulaması yok.
