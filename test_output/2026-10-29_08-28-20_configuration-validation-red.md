# Configuration validation RED

Tarih: 2026-10-29T08:28:20.509893+00:00

Amaç: Geçersiz ayarları reddetme ve sunucu adresini dosyadan okuma. Ortam: Linux, Go 1.24.1.

Komut: `go test ./Lib -run "TestReadConfigurationRejectsInvalidSettings|TestRunUsesConfiguredListeningAddress" -v`

```text
=== RUN   TestReadConfigurationRejectsInvalidSettings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
    configuration_validation_test.go:29: invalid configuration accepted
    configuration_validation_test.go:34: invalid configuration returned partially loaded settings
--- FAIL: TestReadConfigurationRejectsInvalidSettings (0.00s)
=== RUN   TestRunUsesConfiguredListeningAddress
flag provided but not defined: -config
Usage of obs-lan-controller:
  -addr string
    	HTTP dinleme adresi (tasarım sunucusu; varsayılan tüm IPv4 arayüzleri) (default "0.0.0.0:8080")
    configuration_validation_test.go:56: expected bind failure at configured address, got flag provided but not defined: -config
--- FAIL: TestRunUsesConfiguredListeningAddress (0.00s)
FAIL
FAIL	obs-lan-controller/Lib	0.008s
FAIL
```

Çıkış kodu: 1

Gerçek OBS kullanılmadı. Yerel TCP portu test tarafından tutulur.
