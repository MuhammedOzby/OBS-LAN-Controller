# Preview integrated verification
UTC 2026-10-29_16-35-23
Command: `go test ./... -count=1 -timeout=120s`
Exit: 0
```text
ok  	obs-lan-controller	2.094s
ok  	obs-lan-controller/Lib	12.062s
?   	obs-lan-controller/web	[no test files]

```
Normal tests use mocked OBS; preview tests use simulated sender, real WebRTC transport. No real OBS/phone.
