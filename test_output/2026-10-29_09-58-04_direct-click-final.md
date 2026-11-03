# Direct click final
UTC: 2026-10-29T09:58:04.144191+00:00
Linux Go/Node mock tests; no real OBS output commands. Previous timeout root cause: mock DOM remove() did not remove old log children, fixed test double.

`go test ./... -count=1 -timeout=30s` Exit: 0
```
ok  	obs-lan-controller	0.196s
ok  	obs-lan-controller/Lib	0.538s
?   	obs-lan-controller/web	[no test files]

```

`go vet ./...` Exit: 0
```

```
