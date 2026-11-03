# Broadcast backend regression and safety tests

Time: 2026-10-29T09:34:28.040180+00:00
Environment: Linux; CGO enabled for race detector; loopback httptest WebSocket simulation. No real OBS outputs touched.
Purpose: explicit six actions, readback, serialization, stale/no-op prevention, rejection, unsupported compatibility, disconnect no replay, Origin/JSON validation, offline errors; full regression/race/vet.
Command: `CGO_ENABLED=1 GIN_MODE=release go test -race ./... -count=1 -v && go vet ./...`
Exit code: 1

```text
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
FAIL	obs-lan-controller [build failed]
FAIL	obs-lan-controller/Lib [build failed]
?   	obs-lan-controller/web	[no test files]
FAIL
```
Limitations: real OBS output mutations intentionally not verified.
