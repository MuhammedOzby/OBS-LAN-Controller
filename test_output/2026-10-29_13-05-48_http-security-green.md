# HTTP security GREEN

UTC: 2026-10-29T13:05:48.489293+00:00

Purpose: Host guard plus bounded real TCP slow bodies on all action routes.
Environment: Linux Go; synthetic HTTP only, no live OBS/configuration.

Command: `/usr/local/go/bin/go test ./Lib -run TestHTTPServer -count=1 -timeout=30s`
Exit code: 0

```text
ok  	obs-lan-controller/Lib	0.311s

```
Limitations: body tests shorten production 5s ReadTimeout to 100ms; no pairing/authentication claim.
