# Cursor fix and live readback
UTC: 2026-10-29T13:11:41.921964+00:00

## go test ./Lib -run TestDeck -count=1
Exit: 0
```
ok  	obs-lan-controller/Lib	6.394s

```

## OBS_DECK_LIVE=1 go test ./Lib -run ^TestDeckLiveAcceptance$ -v -count=1 -timeout=90s
Exit: 0
```
=== RUN   TestDeckLiveAcceptance
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandler.func1 (2 handlers)
[GIN-debug] GET    /api/deck                 --> obs-lan-controller/Lib.newHandler.deckHandler.func4 (2 handlers)
[GIN-debug] POST   /api/deck/action          --> obs-lan-controller/Lib.newHandler.deckActionHandler.func5 (2 handlers)
[GIN-debug] POST   /api/broadcast            --> obs-lan-controller/Lib.newHandler.broadcastHandler.func6 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandler.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandler.func3 (2 handlers)
    deck_live_test.go:143: live scene verified by HTTP and OBS readback
    deck_live_test.go:154: live visibility verified by HTTP and OBS readback
    deck_live_test.go:154: live visibility verified by HTTP and OBS readback
    deck_live_test.go:157: live mute verified by HTTP and OBS readback
    deck_live_test.go:157: live mute verified by HTTP and OBS readback
    deck_live_test.go:159: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:161: live media-pause verified by HTTP and OBS readback
    deck_live_test.go:162: live media-play verified by HTTP and OBS readback
    deck_live_test.go:164: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:171: live media-stop verified by HTTP and OBS readback
    deck_live_test.go:104: temporary resources absent; initial program scene restored
--- PASS: TestDeckLiveAcceptance (3.49s)
PASS
ok  	obs-lan-controller/Lib	3.493s

```