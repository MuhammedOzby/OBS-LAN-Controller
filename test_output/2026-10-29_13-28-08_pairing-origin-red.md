# Pairing native form RED
UTC: 2026-10-29T13:28:08.326120+00:00
Browser LAN form with invalid dummy code reproduced Forbidden Origin; no OBS mutations or secrets. Failed patch insert made no changes.
go test ./Lib -run TestPairingFormPreservesSameOriginOnPost -count=1
Exit 1
```
--- FAIL: TestPairingFormPreservesSameOriginOnPost (0.00s)
    pairing_origin_test.go:14: form policy = "no-referrer"; want same-origin
FAIL
FAIL	obs-lan-controller/Lib	0.005s
FAIL

```
Unit regression checks Referrer-Policy and untrusted null Origin rejection.