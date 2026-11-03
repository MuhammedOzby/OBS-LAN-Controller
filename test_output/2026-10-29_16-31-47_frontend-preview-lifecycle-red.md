# Preview lifecycle RED
UTC: 2026-10-29T16:31:47.876310+00:00
Command: `/usr/local/go/bin/go test frontend_test.go -run TestFrontendPreviewLifecycle$ -v`
Simulated tests; no OBS actions.
```text
=== RUN   TestFrontendPreviewLifecycle
    frontend_test.go:174: open web/static/preview.js: no such file or directory
--- FAIL: TestFrontendPreviewLifecycle (0.00s)
FAIL
FAIL	command-line-arguments	0.006s
FAIL

```
Exit code: 1
Live OBS/browser video not verified.
