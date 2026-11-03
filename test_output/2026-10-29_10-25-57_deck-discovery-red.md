# Deck discovery RED

UTC: 2026-10-29T10:25:57.767014+00:00

Purpose: new read-only Deck endpoint. Loopback mock websocket only; no live OBS or credentials.

Command: `go test ./Lib -run TestDeckDiscovery -v -count=1 -timeout=30s`

Exit: 1

```text
=== RUN   TestDeckDiscovery
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandler.func1 (2 handlers)
[GIN-debug] POST   /api/broadcast            --> obs-lan-controller/Lib.newHandler.broadcastHandler.func4 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandler.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
    deck_test.go:95: GET = 404 404 page not found
--- FAIL: TestDeckDiscovery (0.00s)
FAIL
FAIL	obs-lan-controller/Lib	0.008s
FAIL

```

Limitations: no live OBS verification.
