# Deck integration verification
UTC: 2026-10-29T13:05:51.811203+00:00
Tests use fixture except explicit live acceptance. Live test restores program and verifies temporary resources absent.

## gofmt -w Lib/deck_live_test.go
Exit: 0
```

```

## go test ./... -count=1
Exit: 1
```
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
--- FAIL: TestBroadcastStructure (0.00s)
    broadcast_frontend_test.go:22: confirmation UI on /deck
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
--- FAIL: TestDeckBrowserRegression (25.05s)
    deck_frontend_test.go:58: Chromium: signal: killed
        [377256:377269:1029/130552.886414:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377269:1029/130552.904373:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377269:1029/130552.928422:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130552.951583:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130552.952063:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130552.953001:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377256:1029/130552.957700:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377256:1029/130552.963502:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        
        (process:377256): GLib-GIO-CRITICAL **: 13:05:52.969: g_settings_schema_source_lookup: assertion 'source != NULL' failed
        [377256:377269:1029/130552.977585:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130553.091410:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.Properties.GetAll: object_path= /org/freedesktop/UPower/devices/DisplayDevice: org.freedesktop.DBus.Error.ServiceUnknown: The name is not activatable
        [377256:377256:1029/130553.094351:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130553.101181:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130553.113357:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130553.113473:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130553.124330:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130553.124934:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130553.135213:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130553.135339:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [377256:377256:1029/130553.137547:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [377256:377269:1029/130555.548571:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
        [377256:377269:1029/130555.559374:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
        [377256:377269:1029/130555.560728:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
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
--- FAIL: TestFrontendPollingAndFailure (0.05s)
    frontend_test.go:56: frontend behavior: exit status 1
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        
        0 !== 1
        
            at [eval]:128:9 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: 0,
          expected: 1,
          operator: 'strictEqual',
          diff: 'simple'
        }
--- FAIL: TestFrontendSceneCommand (0.05s)
    frontend_test.go:87: frontend behavior: exit status 1
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        + actual - expected
        
        + 'undefined'
        - 'function'
        
            at [eval]:127:25 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: 'undefined',
          expected: 'function',
          operator: 'strictEqual',
          diff: 'simple'
        }
--- FAIL: TestFrontendOtherPageNoSnapshot (0.04s)
    frontend_test.go:111: frontend behavior: exit status 1
        TypeError: Cannot read properties of undefined (reading 'url')
            at [eval]:238:42
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
--- FAIL: TestDesignPages (0.00s)
    --- FAIL: TestDesignPages//deck (0.00s)
        main_test.go:42: header should be removed
FAIL
FAIL	obs-lan-controller	25.332s
ok  	obs-lan-controller/Lib	10.941s
?   	obs-lan-controller/web	[no test files]
FAIL

```

## go vet ./...
Exit: 0
```

```

## OBS_DECK_LIVE=1 go test ./Lib -run ^TestDeckLiveAcceptance$ -v -count=1 -timeout=90s
Exit: 1
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
    deck_live_test.go:128: live scene verified by HTTP and OBS readback
    deck_live_test.go:139: live visibility verified by HTTP and OBS readback
    deck_live_test.go:139: live visibility verified by HTTP and OBS readback
    deck_live_test.go:142: live mute verified by HTTP and OBS readback
    deck_live_test.go:142: live mute verified by HTTP and OBS readback
    deck_live_test.go:144: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:146: live media-pause verified by HTTP and OBS readback
    deck_live_test.go:147: live media-play verified by HTTP and OBS readback
    deck_live_test.go:149: live media-restart verified by HTTP and OBS readback
    deck_live_test.go:156: live media-stop verified by HTTP and OBS readback
    deck_live_test.go:68: temporary scene still exists
    deck_live_test.go:81: temporary input still exists
    deck_live_test.go:90: temporary resources removed; initial program scene restored
--- FAIL: TestDeckLiveAcceptance (3.38s)
FAIL
FAIL	obs-lan-controller/Lib	3.380s
FAIL

```