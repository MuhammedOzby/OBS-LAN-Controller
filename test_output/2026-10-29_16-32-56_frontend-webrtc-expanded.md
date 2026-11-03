# Expanded frontend WebRTC verification
UTC: 2026-10-29T16:32:56.807637+00:00
Command: `/usr/local/go/bin/go test frontend_test.go -v`
Simulated Node and HTTP template tests; no OBS actions.
```text
=== RUN   TestFrontendControlStructure
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /api/status               --> obs-lan-controller/Lib.newHandlerWithPreview.func1 (2 handlers)
[GIN-debug] GET    /api/deck                 --> obs-lan-controller/Lib.newHandlerWithPreview.deckHandler.func4 (2 handlers)
[GIN-debug] POST   /api/deck/action          --> obs-lan-controller/Lib.newHandlerWithPreview.deckActionHandler.func5 (2 handlers)
[GIN-debug] POST   /api/broadcast            --> obs-lan-controller/Lib.newHandlerWithPreview.broadcastHandler.func6 (2 handlers)
[GIN-debug] POST   /api/scene                --> obs-lan-controller/Lib.newHandlerWithPreview.func2 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.newHandlerWithPreview.func3 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.newHandlerWithPreview.func3 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.newHandlerWithPreview.func3 (2 handlers)
--- PASS: TestFrontendControlStructure (0.00s)
=== RUN   TestFrontendBrowserServer
    frontend_test.go:48: explicit browser harness only
--- SKIP: TestFrontendBrowserServer (0.00s)
=== RUN   TestFrontendPollingAndFailure
--- PASS: TestFrontendPollingAndFailure (0.04s)
=== RUN   TestFrontendSceneCommand
--- PASS: TestFrontendSceneCommand (0.04s)
=== RUN   TestFrontendOtherPageNoSnapshot
--- PASS: TestFrontendOtherPageNoSnapshot (0.04s)
=== RUN   TestFrontendDOMSelectorsAndDeckGuard
--- PASS: TestFrontendDOMSelectorsAndDeckGuard (0.04s)
=== RUN   TestFrontendPreviewLifecycle
--- PASS: TestFrontendPreviewLifecycle (0.05s)
=== RUN   TestFrontendPreviewStaleAnswer
--- PASS: TestFrontendPreviewStaleAnswer (0.05s)
=== RUN   TestFrontendPreviewBoundedRetry
--- PASS: TestFrontendPreviewBoundedRetry (0.04s)
=== RUN   TestFrontendPreviewICETimeout
--- PASS: TestFrontendPreviewICETimeout (0.04s)
=== RUN   TestFrontendPreviewFetchDeadline
--- PASS: TestFrontendPreviewFetchDeadline (0.05s)
=== RUN   TestFrontendPreviewStatsAndAuth
--- PASS: TestFrontendPreviewStatsAndAuth (0.04s)
PASS
ok  	command-line-arguments	0.441s

```
Exit code: 0
Live OBS/browser video not verified.
