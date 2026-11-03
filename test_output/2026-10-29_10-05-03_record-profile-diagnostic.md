# Real OBS recording profile diagnostic
UTC: 2026-10-29T10:05:03.828562+00:00
Command: OBS_RECORD_DIAGNOSTIC=1 go test ./Lib -run TestRealRecordPauseDiagnostic -v -count=1
Exit: 0
```
=== RUN   TestRealRecordPauseDiagnostic
    record_real_test.go:27: OBS version=33.0.0 WS=5.7.6
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:Simple parameterValue:Simple]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:none parameterValue:none]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:Stream parameterValue:Stream]
    record_real_test.go:27: GetProfileParameter err=<nil> data=map[defaultParameterValue:nvenc parameterValue:nvenc]
    record_real_test.go:27: GetRecordStatus err=<nil> data=map[outputActive:true outputBytes:8.231048e+07 outputDuration:58466 outputPaused:false outputTimecode:00:00:58.466]
    record_real_test.go:27: PauseRecord err=OBS verisi geçersiz data=map[]
    record_real_test.go:29: pause readback 0 err=<nil> data=map[outputActive:true outputBytes:8.269572e+07 outputDuration:58966 outputPaused:false outputTimecode:00:00:58.966]
    record_real_test.go:29: pause readback 1 err=<nil> data=map[outputActive:true outputBytes:8.3080533e+07 outputDuration:59466 outputPaused:false outputTimecode:00:00:59.466]
    record_real_test.go:29: pause readback 2 err=<nil> data=map[outputActive:true outputBytes:8.3465773e+07 outputDuration:59966 outputPaused:false outputTimecode:00:00:59.966]
    record_real_test.go:29: pause readback 3 err=<nil> data=map[outputActive:true outputBytes:8.3850587e+07 outputDuration:60466 outputPaused:false outputTimecode:00:01:00.466]
    record_real_test.go:29: pause readback 4 err=<nil> data=map[outputActive:true outputBytes:8.4235827e+07 outputDuration:60966 outputPaused:false outputTimecode:00:01:00.966]
    record_real_test.go:29: pause readback 5 err=<nil> data=map[outputActive:true outputBytes:8.462064e+07 outputDuration:61466 outputPaused:false outputTimecode:00:01:01.466]
    record_real_test.go:30: resume err=OBS isteği reddetti
--- PASS: TestRealRecordPauseDiagnostic (3.02s)
PASS
ok  	obs-lan-controller/Lib	3.020s

```
