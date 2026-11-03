# Real OBS recording profile reload cycle
UTC 2026-10-29T10:07:59.504213+00:00
HQ+x264 profile configuration reloaded via temporary profile switch (temporary profile removed). User authorized actual OBS operations. Streaming was inactive and not started.
OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run ^TestRealRecordPauseResume$ -v -count=1 -timeout=50s
Exit 1
=== RUN   TestRealRecordPauseResume
    record_pause_real_test.go:24: Recording quality verified HQ; streaming settings unchanged
    record_pause_real_test.go:28: Recording encoder verified x264 (previous nvenc did not start)
    record_pause_real_test.go:35: Original profile reselected to reload recording output configuration
    record_pause_real_test.go:39: start-record HTTP=200 active=true paused=false duration=0 bytes=0
    record_pause_real_test.go:39: pause-record HTTP=502 active=false paused=false duration=0 bytes=0
    record_pause_real_test.go:40: pause-record not verified
--- FAIL: TestRealRecordPauseResume (5.23s)
FAIL
FAIL	obs-lan-controller/Lib	5.242s
FAIL
