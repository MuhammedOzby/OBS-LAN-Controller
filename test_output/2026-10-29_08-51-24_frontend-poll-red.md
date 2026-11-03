# frontend-poll-red

UTC: 2026-10-29T08:51:24.759522+00:00

Komut: /usr/local/go/bin/go test . -run TestFrontendPollingAndFailure -v

Çıkış: 1

```text
=== RUN   TestFrontendPollingAndFailure
    frontend_test.go:37: frontend behavior: exit status 1
        AssertionError [ERR_ASSERTION]: Expected values to be strictly equal:
        
        0 !== 1
        
            at [eval]:19:9 {
          generatedMessage: true,
          code: 'ERR_ASSERTION',
          actual: 0,
          expected: 1,
          operator: 'strictEqual',
          diff: 'simple'
        }
--- FAIL: TestFrontendPollingAndFailure (0.05s)
FAIL
FAIL	obs-lan-controller	0.051s
FAIL
```

Sınır: Simülasyon: gerçek OBS yok; geometri ayrı doğrulanır.
