# Restore original recording settings
UTC 2026-10-29T10:09:17.147068+00:00
Restore original RecQuality Stream and RecEncoder nvenc; stop test recording and read settings/status. Streaming untouched. Diagnostic pause/resume while inactive expected rejection.
OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run ^TestRealRecordPauseDiagnostic$ -v -count=1
Exit 0
=== RUN   TestRealRecordPauseDiagnostic
    record_real_test.go:31: OBS version=33.0.0 WS=5.7.6
    record_real_test.go:31: GetProfileParameter err=<nil> data=map[defaultParameterValue:Simple parameterValue:Simple]
    record_real_test.go:31: GetProfileParameter err=<nil> data=map[defaultParameterValue:none parameterValue:none]
    record_real_test.go:31: GetProfileParameter err=<nil> data=map[defaultParameterValue:Stream parameterValue:Stream]
    record_real_test.go:31: GetProfileParameter err=<nil> data=map[defaultParameterValue:nvenc parameterValue:nvenc]
    record_real_test.go:31: GetRecordStatus err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:31: PauseRecord err=OBS verisi geçersiz data=map[]
    record_real_test.go:33: pause readback 0 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:33: pause readback 1 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:33: pause readback 2 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:33: pause readback 3 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:33: pause readback 4 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:33: pause readback 5 err=<nil> data=map[outputActive:false outputBytes:3.3494125e+07 outputDuration:0 outputPaused:false outputTimecode:00:00:00.000]
    record_real_test.go:34: resume err=OBS isteği reddetti
--- PASS: TestRealRecordPauseDiagnostic (3.02s)
PASS
ok  	obs-lan-controller/Lib	3.027s
