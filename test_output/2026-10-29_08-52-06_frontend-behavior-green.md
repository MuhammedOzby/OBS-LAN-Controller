# frontend-behavior-green

UTC: 2026-10-29T08:52:06.273128+00:00

Komut: /usr/local/go/bin/go test . -run TestFrontend -v

Çıkış: 0

```text
=== RUN   TestFrontendControlStructure
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandler.func1 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandler.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
--- PASS: TestFrontendControlStructure (0.00s)
=== RUN   TestFrontendBrowserServer
    frontend_test.go:31: explicit browser harness only
--- SKIP: TestFrontendBrowserServer (0.00s)
=== RUN   TestFrontendPollingAndFailure
--- PASS: TestFrontendPollingAndFailure (0.04s)
=== RUN   TestFrontendSceneCommand
--- PASS: TestFrontendSceneCommand (0.04s)
=== RUN   TestFrontendOtherPageNoSnapshot
--- PASS: TestFrontendOtherPageNoSnapshot (0.04s)
PASS
ok  	obs-lan-controller	0.120s
```

Sınır: Simülasyon: gerçek OBS yok; geometri ayrı doğrulanır.
