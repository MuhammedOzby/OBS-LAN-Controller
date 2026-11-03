# Deck discovery GREEN

UTC: 2026-10-29T10:26:27.638511+00:00

Purpose: read-only Deck discovery. Loopback mock websocket only; no live OBS.

Command: `go test ./Lib -run TestDeckDiscovery -v -count=1 -timeout=30s`

Exit: 0

```text
=== RUN   TestDeckDiscovery
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandler.func1 (2 handlers)
[GIN-debug] GET    /api/deck                 --> obs-lan-controller/Lib.newHandler.deckHandler.func4 (2 handlers)
[GIN-debug] POST   /api/broadcast            --> obs-lan-controller/Lib.newHandler.broadcastHandler.func5 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandler.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
--- PASS: TestDeckDiscovery (0.00s)
PASS
ok  	obs-lan-controller/Lib	0.007s

```

Limitations: no live OBS verification.
