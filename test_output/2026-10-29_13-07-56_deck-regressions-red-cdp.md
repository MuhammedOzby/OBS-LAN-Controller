# Deck regression RED

UTC: 2026-10-29T13:07:56.971480+00:00

Purpose: reproduce portrait overlap, missing target auto-selection, stale identity using real Chromium CDP and model test. No live OBS/configuration reads.

Command: `/usr/local/go/bin/go test . -run TestDeck(BrowserRegression|Model)$ -v -count=1`

Exit: 1

```text
=== RUN   TestDeckModel
    deck_frontend_test.go:49: exit status 1
        node:internal/assert/utils:146
          throw error;
          ^
        
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        + actual - expected
        
        + undefined
        - 'identity-logo'
        
            at [eval]:4:9
            at runScriptInThisContext (node:internal/vm:219:10)
            at node:internal/process/execution:485:12
            at [eval]-wrapper:6:24
            at runScriptInContext (node:internal/process/execution:483:60)
            at evalFunction (node:internal/process/execution:317:30)
            at evalTypeScript (node:internal/process/execution:329:3)
            at node:internal/main/eval_string:71:3 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: undefined,
          expected: 'identity-logo',
          operator: 'strictEqual',
          diff: 'simple'
        }
        
        Node.js v26.12.1
--- FAIL: TestDeckModel (0.04s)
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
    deck_frontend_test.go:85: browser regression: FAIL portrait: last key overflows grid: {"grid":{"x":29,"y":204.921875,"width":332,"height":181.515625,"top":204.921875,"right":361,"bottom":386.4375,"left":29},"last":{"x":255.65625,"y":552.921875,"width":105.34375,"height":108,"top":552.921875,"right":361,"bottom":660.921875,"left":255.65625}} | missing scene editor: silently selected first remaining scene | stale visibility identity editor: rebound reused scene item id
--- FAIL: TestDeckBrowserRegression (0.59s)
FAIL
FAIL	obs-lan-controller	0.637s
FAIL

```

Limits: fixtures, not live OBS.
