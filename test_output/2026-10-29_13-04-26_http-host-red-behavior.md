# HTTP trusted Host behavioral RED

UTC: 2026-10-29T13:04:26.267391+00:00

Purpose: show rejection assertions fail with existing unguarded behavior (constructor extracted).
Environment: Linux, Go; synthetic HTTP only; no live OBS/configuration.

Command: `/usr/local/go/bin/go test ./Lib -run TestHTTPServer -count=1`
Exit code: 1

```text
--- FAIL: TestHTTPServerTrustedHosts (0.00s)
    --- FAIL: TestHTTPServerTrustedHosts/attacker.example:8080 (0.00s)
        http_security_test.go:45: Host "attacker.example:8080": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/controller.test (0.00s)
        http_security_test.go:45: Host "controller.test": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/localhost.attacker.example (0.00s)
        http_security_test.go:45: Host "localhost.attacker.example": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/192.0.2.123:8080 (0.00s)
        http_security_test.go:45: Host "192.0.2.123:8080": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/[2001:db8::123]:8080 (0.00s)
        http_security_test.go:45: Host "[2001:db8::123]:8080": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/0.0.0.0:8080 (0.00s)
        http_security_test.go:45: Host "0.0.0.0:8080": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/[::]:8080 (0.00s)
        http_security_test.go:45: Host "[::]:8080": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/127.0.0.1:invalid (0.00s)
        http_security_test.go:45: Host "127.0.0.1:invalid": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/localhost:65536 (0.00s)
        http_security_test.go:45: Host "localhost:65536": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/localhost: (0.00s)
        http_security_test.go:45: Host "localhost:": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/localhost@attacker.example (0.00s)
        http_security_test.go:45: Host "localhost@attacker.example": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/#00 (0.00s)
        http_security_test.go:45: Host "": status=204 called=true, want 403
    --- FAIL: TestHTTPServerTrustedHosts/::1 (0.00s)
        http_security_test.go:45: Host "::1": status=204 called=true, want 403
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
--- FAIL: TestHTTPServerGuardPreservesOriginChecksAndFixtures (0.00s)
    http_security_test.go:90: /api/scene Host="attacker.example" Origin="http://attacker.example": 400, want 403
    http_security_test.go:90: /api/broadcast Host="attacker.example" Origin="http://attacker.example": 400, want 403
    http_security_test.go:90: /api/deck/action Host="attacker.example" Origin="http://attacker.example": 400, want 403
FAIL
FAIL	obs-lan-controller/Lib	0.006s
FAIL

```
Limitations: pairing/authentication unimplemented. Initial red also recorded missing constructor compile failure.
