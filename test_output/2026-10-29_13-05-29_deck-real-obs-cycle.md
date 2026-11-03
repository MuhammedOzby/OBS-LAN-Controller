# Real demo OBS Deck acceptance
UTC: 2026-10-29T13:05:29.277994+00:00
User explicitly authorized all demo OBS use. Temporary scene/media input, no stream/record mutations.

## ffmpeg -hide_banner -loglevel error -f lavfi -i color=c=blue:s=320x180:r=25 -f lavfi -i sine=frequency=440:sample_rate=44100 -t 30 -c:v mpeg4 -c:a aac -y test_output/deck-acceptance.mp4
Exit: 0
```

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
    deck_live_test.go:102: live scene verified by HTTP and OBS readback
    deck_live_test.go:113: live visibility verified by HTTP and OBS readback
    deck_live_test.go:113: live visibility verified by HTTP and OBS readback
    deck_live_test.go:116: live mute verified by HTTP and OBS readback
    deck_live_test.go:116: live mute verified by HTTP and OBS readback
    deck_live_test.go:118: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:120: live media-pause verified by HTTP and OBS readback
    deck_live_test.go:121: live media-play verified by HTTP and OBS readback
    deck_live_test.go:123: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:130: live media-stop verified by HTTP and OBS readback
    deck_live_test.go:64: temporary resources removed; initial program scene restored
--- PASS: TestDeckLiveAcceptance (3.48s)
PASS
ok  	obs-lan-controller/Lib	3.484s

```