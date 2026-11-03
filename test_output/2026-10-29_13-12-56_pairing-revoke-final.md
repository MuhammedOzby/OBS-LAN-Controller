# Pairing and local revocation final
UTC: 2026-10-29T13:12:56.614092+00:00
Purpose: pairing, socket-peer trust, expiry, finite storage, global failures, Origin/body limits, local revoke, downstream guards and regression. Synthetic HTTP/mock OBS only. No live OBS/credentials.

## `GIN_MODE=release go test ./Lib -run TestPairing -count=1 -v`
Exit: 0
```text
=== RUN   TestPairingLifecycle
--- PASS: TestPairingLifecycle (0.00s)
=== RUN   TestPairingExpiryRotationAndLimits
--- PASS: TestPairingExpiryRotationAndLimits (0.01s)
=== RUN   TestPairingOriginsBodyAndTLS
--- PASS: TestPairingOriginsBodyAndTLS (0.00s)
=== RUN   TestPairingConcurrentSingleUse
--- PASS: TestPairingConcurrentSingleUse (0.00s)
=== RUN   TestPairingPreservesDownstreamOriginAndHostGuard
--- PASS: TestPairingPreservesDownstreamOriginAndHostGuard (0.00s)
=== RUN   TestPairingLocalRevokeAll
--- PASS: TestPairingLocalRevokeAll (0.00s)
=== RUN   TestPairingRemoteAccessBoundary
--- PASS: TestPairingRemoteAccessBoundary (0.00s)
PASS
ok  	obs-lan-controller/Lib	0.014s

```

## `GIN_MODE=release go test ./... -count=1`
Exit: 0
```text
ok  	obs-lan-controller	1.278s
ok  	obs-lan-controller/Lib	10.939s
?   	obs-lan-controller/web	[no test files]

```

## `go vet ./...`
Exit: 0
```text

```

## `go build ./...`
Exit: 0
```text

```

Limitations: race blocked (CGO initially disabled; retry lacked gcc/cc/clang/zig). Physical LAN, browser appearance, trusted HTTPS and real OBS remain parent acceptance. HTTP is unencrypted. Sessions last 12h, max 128, global failed attempt limit 20/min, 5-minute single-use code; local same-origin revoke clears sessions/pending code without resetting failure throttle.
