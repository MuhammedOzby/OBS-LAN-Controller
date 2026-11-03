# Pairing lifecycle RED
UTC: 2026-10-29T13:09:51.929437+00:00
Purpose: synthetic pairing and authenticated readback; no OBS.
Command: `go test ./Lib -run TestPairingLifecycle -count=1`
Exit: 1
```text
--- FAIL: TestPairingLifecycle (0.00s)
    pairing_test.go:35: local pairing code missing
FAIL
FAIL	obs-lan-controller/Lib	0.004s
FAIL

```
Limitations: no live OBS/browser.
