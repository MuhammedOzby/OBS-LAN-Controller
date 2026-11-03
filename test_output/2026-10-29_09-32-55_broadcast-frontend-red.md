# Broadcast frontend RED

UTC: 2026-10-29T09:32:55.304713+00:00
Amaç: Eksik yayın yönetimi davranışlarını kanıtlamak. Ortam: Linux, Go/Node; simülasyon, gerçek OBS komutu yok.
Komut: `go test . -run TestBroadcast -count=1`
Çıkış: 1
```text
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
--- FAIL: TestBroadcastStructure (0.00s)
    broadcast_frontend_test.go:17: /yayin-yonetimi marker id="broadcast-dashboard"
    broadcast_frontend_test.go:17: /yayin-yonetimi marker id="broadcast-confirm"
    broadcast_frontend_test.go:17: /yayin-yonetimi marker data-broadcast-action="start-stream"
    broadcast_frontend_test.go:17: /yayin-yonetimi marker data-broadcast-action="stop-record"
    broadcast_frontend_test.go:17: /yayin-yonetimi marker id="broadcast-log"
--- FAIL: TestBroadcastFrontendBehavior (0.00s)
    broadcast_frontend_test.go:23: open web/static/broadcast.js: no such file or directory
FAIL
FAIL	obs-lan-controller	0.005s
FAIL

```
