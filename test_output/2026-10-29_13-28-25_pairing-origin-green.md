# Pairing Origin fix verification
UTC: 2026-10-29T13:28:25.987944+00:00
Unit/fixture/browser suite, no live OBS mutation.

go test ./... -count=1 -timeout=120s
Exit 0
```
ok  	obs-lan-controller	1.198s
ok  	obs-lan-controller/Lib	10.932s
?   	obs-lan-controller/web	[no test files]

```

go vet ./...
Exit 0
```

```

go build -o test_output/obs-lan-controller-demo .
Exit 0
```

```