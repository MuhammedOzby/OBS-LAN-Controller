# HTTP security race

UTC: 2026-10-29T13:07:10.984437+00:00
Environment: Linux Go1.24.1 CGO_ENABLED=1; OBS_DECK_LIVE disabled.
Purpose: race-check shared HTTP security tests. No live OBS or repository configuration.
Command: `go test -race ./Lib -run TestHTTPServer -count=1 -timeout=60s`
Exit code: 1

```text
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
FAIL	obs-lan-controller/Lib [build failed]
FAIL

```
Limitations: no pairing/authentication, startup interface snapshot; full frontend regression failures recorded separately and left to frontend owner.
