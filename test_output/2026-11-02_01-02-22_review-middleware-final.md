# Final middleware and browser verification

Date: 2026-11-02T01:02:36+03:00

Purpose: full Go/Node/Chromium tests, production wrapper security checks, vet and build.

Commands: `go test ./... -count=1` exit 0; `go vet ./...` exit 0; `go build -o /tmp/opencode/obs-controller-review .` exit 0. Node and CHROME_BIN paths supplied.

Mock OBS and browser fixtures only; no real OBS connection. Physical devices and real OBS remain unverified.
