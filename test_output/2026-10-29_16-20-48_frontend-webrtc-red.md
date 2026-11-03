# Frontend WebRTC RED

UTC: 2026-10-29T16:20:48.525702+00:00

Command: `/usr/local/go/bin/go test frontend_test.go -run TestFrontend(ControlStructure|PollingAndFailure)$ -v`

Simulated Node/HTTP template tests only; no OBS actions.

```text
=== RUN   TestFrontendControlStructure
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
    frontend_test.go:26: control missing "id=\"program-video\""
    frontend_test.go:26: control missing "autoplay muted playsinline"
    frontend_test.go:26: control missing "id=\"preview-status\""
    frontend_test.go:26: control missing "id=\"preview-fps\""
    frontend_test.go:26: control missing "/static/preview.js"
--- FAIL: TestFrontendControlStructure (0.00s)
=== RUN   TestFrontendPollingAndFailure
    frontend_test.go:58: frontend behavior: exit status 1
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        + actual - expected

        + '/api/status?preview=1'
        - '/api/status?preview=0'
                               ^

            at [eval]:130:42 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: '/api/status?preview=1',
          expected: '/api/status?preview=0',
          operator: 'strictEqual',
          diff: 'simple'
        }
--- FAIL: TestFrontendPollingAndFailure (0.05s)
FAIL
FAIL	command-line-arguments	0.055s
FAIL

```
Exit code: 1
