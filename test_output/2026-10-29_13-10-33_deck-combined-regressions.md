# Deck combined regressions
UTC: 2026-10-29T13:10:33.677690+00:00
Fixture and headless Chromium; no live OBS.

## go test ./... -count=1 -timeout=120s
Exit: 0
```
ok  	obs-lan-controller	1.231s
ok  	obs-lan-controller/Lib	10.931s
?   	obs-lan-controller/web	[no test files]

```

## go vet ./...
Exit: 0
```

```

## go build -o test_output/obs-lan-controller-demo .
Exit: 0
```

```