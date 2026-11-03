# Review final suite including responsive dialogs

Date: 2026-11-02T01:01:52+03:00

Purpose: full Go/Node/headless Chromium suite including portrait and short landscape dialogs; vet and build.

Commands: `go test ./... -count=1` exit 0; `go vet ./...` exit 0; `go build -o /tmp/opencode/obs-controller-review .` exit 0. Node and CHROME_BIN paths supplied.

Mock OBS and browser fixtures only; no real OBS connection. Physical devices and real OBS unverified.
