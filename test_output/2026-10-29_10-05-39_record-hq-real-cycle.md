# Real OBS HQ recording cycle
UTC: 2026-10-29T10:05:39.755517+00:00
User authorized actual OBS changes. Stop test recording, change SimpleOutput RecQuality Stream to HQ, read back, run start/pause/resume/stop. Stream settings untouched.
Command: OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run ^TestRealRecordPauseResume$ -v -count=1 -timeout=50s
Exit: 1
```
=== RUN   TestRealRecordPauseResume
    record_pause_real_test.go:20: stop initial recording HTTP 502
--- FAIL: TestRealRecordPauseResume (0.41s)
FAIL
FAIL	obs-lan-controller/Lib	0.414s
FAIL

```
