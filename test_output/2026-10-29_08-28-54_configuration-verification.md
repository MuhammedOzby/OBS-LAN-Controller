# Yapılandırma doğrulaması

Tarih: 2026-10-29T08:28:54.311959+00:00

Amaç: Dosya okuma, hatalı ayarlar, dosyadan dinleme adresi ve regresyon. Ortam: Linux, Go 1.24.1.

## `go test ./... -v`

```text
=== RUN   TestHandlerUsesGin
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
--- PASS: TestHandlerUsesGin (0.00s)
=== RUN   TestDesignPages
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
=== RUN   TestDesignPages//
=== RUN   TestDesignPages//deck
=== RUN   TestDesignPages//yayin-yonetimi
--- PASS: TestDesignPages (0.00s)
    --- PASS: TestDesignPages// (0.00s)
    --- PASS: TestDesignPages//deck (0.00s)
    --- PASS: TestDesignPages//yayin-yonetimi (0.00s)
PASS
ok  	obs-lan-controller	0.007s
=== RUN   TestReadConfiguration
--- PASS: TestReadConfiguration (0.00s)
=== RUN   TestReadConfigurationRejectsInvalidSettings
--- PASS: TestReadConfigurationRejectsInvalidSettings (0.00s)
=== RUN   TestRunUsesConfiguredListeningAddress
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
2026/10/29 08:28:55 Tasarım sunucusu: http://127.0.0.1:39481 — OBS bağlantısı yok
--- PASS: TestRunUsesConfiguredListeningAddress (0.00s)
PASS
ok  	obs-lan-controller/Lib	0.008s
?   	obs-lan-controller/web	[no test files]
```

Çıkış kodu: 0

## `go test -race ./...`

```text
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
```

Çıkış kodu: 2

## `go vet ./...`

```text
```

Çıkış kodu: 0

## `go build -o obs-lan-controller .`

```text
```

Çıkış kodu: 0

## `git check-ignore configuration.json`

```text
configuration.json
```

Çıkış kodu: 0

## `git diff --check`

```text
```

Çıkış kodu: 0

Gerçek OBS kullanılmadı. Testler dosya ve yerel TCP kullanır; OBS bağlantısı uygulanmadı. Git kontrolleri uygulama testi değildir.
