# Broadcast backend regression and safety tests

Time: 2026-10-29T09:34:23.242846+00:00
Environment: Linux; loopback httptest WebSocket simulation. No real OBS outputs touched.
Purpose: explicit six actions, readback, serialization, stale/no-op prevention, rejection, unsupported status compatibility, disconnect no replay, Origin and JSON validation, offline errors; full regression race detector and vet.
Command: `GIN_MODE=release go test -race ./... -count=1 -v && go vet ./...`
Exit code: 2

```text
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
```
Limitations: real OBS output mutations intentionally not verified.
