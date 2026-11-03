# Deck browser retry

UTC: 2026-10-29T13:06:07.778067+00:00

Fixture only; no live OBS/configuration. DOMContentLoaded checks and bounded Chromium dump.

Command: `/usr/local/go/bin/go test . -run TestDeckBrowserRegression$ -v -count=1`

Exit: 1

```text
=== RUN   TestDeckBrowserRegression
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
    deck_frontend_test.go:58: Chromium: signal: killed
        [376882:376895:1029/130542.815602:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376895:1029/130542.822876:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376895:1029/130542.850673:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130542.872725:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130542.872852:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130542.873154:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376882:1029/130542.880523:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        
        (process:376882): GLib-GIO-CRITICAL **: 13:05:42.889: g_settings_schema_source_lookup: assertion 'source != NULL' failed
        [376882:376895:1029/130542.896319:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130543.011431:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376882:1029/130543.016491:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130543.016593:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130543.029170:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.Properties.GetAll: object_path= /org/freedesktop/UPower/devices/DisplayDevice: org.freedesktop.DBus.Error.ServiceUnknown: The name is not activatable
        [376882:376882:1029/130543.032866:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130543.033043:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130543.044309:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130543.044461:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130543.056143:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130543.056243:ERROR:dbus/bus.cc:405] Failed to connect to the bus: Could not parse server address: Unknown address type (examples of valid types are "tcp" and on UNIX "unix")
        [376882:376882:1029/130543.065247:ERROR:dbus/object_proxy.cc:572] Failed to call method: org.freedesktop.DBus.NameHasOwner: object_path= /org/freedesktop/DBus: unknown error type: 
        [376882:376895:1029/130545.330617:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
        [376882:376895:1029/130545.355722:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
        [376882:376895:1029/130545.360347:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: PHONE_REGISTRATION_ERROR
        [376882:376895:1029/130607.554092:ERROR:google_apis/gcm/engine/registration_request.cc:290] Registration response error message: DEPRECATED_ENDPOINT
--- FAIL: TestDeckBrowserRegression (25.03s)
FAIL
FAIL	obs-lan-controller	25.040s
FAIL

```

Limitations: fixture browser, not real OBS.
