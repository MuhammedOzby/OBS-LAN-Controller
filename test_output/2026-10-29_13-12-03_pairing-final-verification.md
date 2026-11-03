# Pairing final verification
UTC: 2026-10-29T13:12:03.827146+00:00
Purpose: retry race with cgo, regression rerun during concurrent parent Deck work. Synthetic HTTP/mock OBS; no credentials/live OBS.

## `CGO_ENABLED=1 go test -race ./Lib -run "TestPairing|TestHTTPServer" -count=1`
Exit: 1
```text
# runtime/cgo
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in $PATH
FAIL	obs-lan-controller/Lib [build failed]
FAIL

```

## `GIN_MODE=release go test ./... -count=1`
Exit: 0
```text
ok  	obs-lan-controller	1.449s
ok  	obs-lan-controller/Lib	10.965s
?   	obs-lan-controller/web	[no test files]

```

Limitations: physical LAN/mobile/browser and real OBS not exercised. Trusted HTTPS provisioning, per-session revocation, QR and persistence out of scope. HTTP is unencrypted.
