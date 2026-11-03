# frontend-scene-red

UTC: 2026-10-29T08:51:54.375293+00:00

Komut: /usr/local/go/bin/go test . -run TestFrontendSceneCommand -v

Çıkış: 1

```text
=== RUN   TestFrontendSceneCommand
    frontend_test.go:62: frontend behavior: exit status 1
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        + actual - expected
        
        + 'undefined'
        - 'function'
        
            at [eval]:93:25 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: 'undefined',
          expected: 'function',
          operator: 'strictEqual',
          diff: 'simple'
        }
--- FAIL: TestFrontendSceneCommand (0.04s)
FAIL
FAIL	obs-lan-controller	0.049s
FAIL
```

Sınır: Simülasyon: gerçek OBS yok; geometri ayrı doğrulanır.
