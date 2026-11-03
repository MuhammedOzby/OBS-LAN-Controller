# Real OBS diagnostic
UTC: 2026-10-29T10:04:44.820490+00:00
User authorized recording mutations, stream untouched. Command: go test ./Lib -run TestRealRecordPauseDiagnostic -v -count=1
Exit: 0
```
=== RUN   TestRealRecordPauseDiagnostic
    record_real_test.go:25: OBS version=33.0.0 WS=5.7.6
    record_real_test.go:25: GetProfileParameter err=<nil> data=map[defaultParameterValue:Simple parameterValue:Simple]
    record_real_test.go:25: GetProfileParameter err=<nil> data=map[defaultParameterValue:none parameterValue:none]
    record_real_test.go:25: GetRecordStatus err=<nil> data=map[outputActive:true outputBytes:6.8091413e+07 outputDuration:39999 outputPaused:false outputTimecode:00:00:39.999]
    record_real_test.go:25: PauseRecord err=OBS verisi geçersiz data=map[]
    record_real_test.go:27: pause readback 0 err=<nil> data=map[outputActive:true outputBytes:6.8476227e+07 outputDuration:40499 outputPaused:false outputTimecode:00:00:40.499]
    record_real_test.go:27: pause readback 1 err=<nil> data=map[outputActive:true outputBytes:6.8861467e+07 outputDuration:40999 outputPaused:false outputTimecode:00:00:40.999]
    record_real_test.go:27: pause readback 2 err=<nil> data=map[outputActive:true outputBytes:6.924628e+07 outputDuration:41499 outputPaused:false outputTimecode:00:00:41.499]
    record_real_test.go:27: pause readback 3 err=<nil> data=map[outputActive:true outputBytes:6.963152e+07 outputDuration:41999 outputPaused:false outputTimecode:00:00:41.999]
    record_real_test.go:27: pause readback 4 err=<nil> data=map[outputActive:true outputBytes:7.0016333e+07 outputDuration:42499 outputPaused:false outputTimecode:00:00:42.499]
    record_real_test.go:27: pause readback 5 err=<nil> data=map[outputActive:true outputBytes:7.0426573e+07 outputDuration:43033 outputPaused:false outputTimecode:00:00:43.033]
    record_real_test.go:28: resume err=OBS isteği reddetti
--- PASS: TestRealRecordPauseDiagnostic (3.02s)
PASS
ok  	obs-lan-controller/Lib	3.023s

```
