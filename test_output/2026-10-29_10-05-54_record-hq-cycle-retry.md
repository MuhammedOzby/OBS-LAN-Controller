# Real OBS HQ cycle retry
UTC: 2026-10-29T10:05:54.033432+00:00
Previous stop completed after 400ms readback window; GET confirmed inactive before retry. OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run ^TestRealRecordPauseResume$ -v -count=1 -timeout=50s
Exit: 1
```
=== RUN   TestRealRecordPauseResume
    record_pause_real_test.go:24: Recording quality verified HQ; streaming settings unchanged
    record_pause_real_test.go:28: start-record HTTP=502 active=false paused=false duration=0 bytes=0
    record_pause_real_test.go:29: start-record not verified
--- FAIL: TestRealRecordPauseResume (1.42s)
FAIL
FAIL	obs-lan-controller/Lib	1.420s
FAIL

```
