# Preview integrated verification
UTC 2026-10-29_16-36-04
Command: `go test -race ./Lib -run TestPreview -count=1 -timeout=90s`
Exit: 0
```text
ok  	obs-lan-controller/Lib	2.180s

```
Normal tests use mocked OBS; preview tests use simulated sender, real WebRTC transport. No real OBS/phone.
