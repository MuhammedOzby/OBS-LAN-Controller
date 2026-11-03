# Post reload real diagnostic
UTC 2026-10-29T10:08:14.939418+00:00
Explicit opt-in Go real OBS diagnostic. Exit 0
=== RUN   TestRealRecordPauseDiagnostic
    record_real_test.go:27: OBS version=33.0.0 WS=5.7.6
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:Simple parameterValue:Simple]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:none parameterValue:none]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:Stream parameterValue:Stream]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:nvenc parameterValue:nvenc]
    record_real_test.go:27: GetRecordStatus err=<nil> data=map[outputActive:true outputBytes:1.0342951e+07 outputDuration:13433 outputPaused:false outputTimecode:00:00:13.433]
    record_real_test.go:27: PauseRecord err=OBS verisi geçersiz data=map[]
    record_real_test.go:29: pause readback 0 err=<nil> data=map[outputActive:true outputBytes:1.0727765e+07 outputDuration:13933 outputPaused:false outputTimecode:00:00:13.933]
    record_real_test.go:29: pause readback 1 err=<nil> data=map[outputActive:true outputBytes:1.1113005e+07 outputDuration:14433 outputPaused:false outputTimecode:00:00:14.433]
    record_real_test.go:29: pause readback 2 err=<nil> data=map[outputActive:true outputBytes:1.1497818e+07 outputDuration:14933 outputPaused:false outputTimecode:00:00:14.933]
    record_real_test.go:29: pause readback 3 err=<nil> data=map[outputActive:true outputBytes:1.1883058e+07 outputDuration:15433 outputPaused:false outputTimecode:00:00:15.433]
    record_real_test.go:29: pause readback 4 err=<nil> data=map[outputActive:true outputBytes:1.2267871e+07 outputDuration:15933 outputPaused:false outputTimecode:00:00:15.933]
    record_real_test.go:29: pause readback 5 err=<nil> data=map[outputActive:true outputBytes:1.2653111e+07 outputDuration:16433 outputPaused:false outputTimecode:00:00:16.433]
    record_real_test.go:30: resume err=OBS isteği reddetti
--- PASS: TestRealRecordPauseDiagnostic (3.02s)
PASS
ok  	obs-lan-controller/Lib	3.021s
