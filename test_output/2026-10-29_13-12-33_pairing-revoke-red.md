# Pairing revoke RED
UTC: 2026-10-29T13:12:33.157849+00:00
Purpose: local-only same-origin revoke regression. Synthetic HTTP; no OBS.
Command: `GIN_MODE=release go test ./Lib -run TestPairingLocalRevokeAll -count=1`
Exit: 1
```text
--- FAIL: TestPairingLocalRevokeAll (0.00s)
    pairing_test.go:257: remote revoke got 401
FAIL
FAIL	obs-lan-controller/Lib	0.005s
FAIL

```
Limitations: no live browser/LAN/OBS.
