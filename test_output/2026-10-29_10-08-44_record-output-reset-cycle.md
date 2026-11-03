# Real OBS unchanged resolution output reset
UTC 2026-10-29T10:08:44.207850+00:00
User authorized recording tests. Set HQ/x264, reset video at same output resolution, then cycle. Stream remains off.
OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run ^TestRealRecordPauseResume$ -v -count=1
Exit 1
=== RUN   TestRealRecordPauseResume
    record_pause_real_test.go:24: Recording quality verified HQ; streaming settings unchanged
    record_pause_real_test.go:28: Recording encoder verified x264 (previous nvenc did not start)
    record_pause_real_test.go:32: Video reset at unchanged resolution to apply output configuration
    record_pause_real_test.go:36: start-record HTTP=502 active=false paused=false duration=0 bytes=0
    record_pause_real_test.go:37: start-record not verified
--- FAIL: TestRealRecordPauseResume (4.90s)
FAIL
FAIL	obs-lan-controller/Lib	4.908s
FAIL
