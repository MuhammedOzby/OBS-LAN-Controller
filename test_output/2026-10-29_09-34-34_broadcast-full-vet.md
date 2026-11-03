# Broadcast backend full tests and vet

Time: 2026-10-29T09:34:34.508824+00:00
Environment: Linux; loopback httptest WS simulation only. No real OBS outputs touched.
Purpose: six explicit actions, readback, concurrency, no-op/stale, rejection, unsupported compatibility, disconnect no replay, Origin/confirmation/whitelist, regression.
Command: `GIN_MODE=release go test ./... -count=1 -v && go vet ./...`
Exit code: 0

```text
=== RUN   TestBroadcastStructure
--- PASS: TestBroadcastStructure (0.00s)
=== RUN   TestBroadcastFrontendBehavior
--- PASS: TestBroadcastFrontendBehavior (0.07s)
=== RUN   TestFrontendControlStructure
--- PASS: TestFrontendControlStructure (0.00s)
=== RUN   TestFrontendBrowserServer
    frontend_test.go:46: explicit browser harness only
--- SKIP: TestFrontendBrowserServer (0.00s)
=== RUN   TestFrontendPollingAndFailure
--- PASS: TestFrontendPollingAndFailure (0.05s)
=== RUN   TestFrontendSceneCommand
--- PASS: TestFrontendSceneCommand (0.05s)
=== RUN   TestFrontendOtherPageNoSnapshot
--- PASS: TestFrontendOtherPageNoSnapshot (0.04s)
=== RUN   TestHandlerUsesGin
--- PASS: TestHandlerUsesGin (0.00s)
=== RUN   TestDesignPages
=== RUN   TestDesignPages//
=== RUN   TestDesignPages//deck
=== RUN   TestDesignPages//yayin-yonetimi
--- PASS: TestDesignPages (0.00s)
    --- PASS: TestDesignPages// (0.00s)
    --- PASS: TestDesignPages//deck (0.00s)
    --- PASS: TestDesignPages//yayin-yonetimi (0.00s)
PASS
ok  	obs-lan-controller	0.217s
=== RUN   TestBroadcastSnapshotRecording
--- PASS: TestBroadcastSnapshotRecording (0.00s)
=== RUN   TestBroadcastStartReadback
--- PASS: TestBroadcastStartReadback (0.10s)
=== RUN   TestBroadcastActionsAndNoOps
=== RUN   TestBroadcastActionsAndNoOps/start-stream
=== RUN   TestBroadcastActionsAndNoOps/stop-stream
=== RUN   TestBroadcastActionsAndNoOps/start-record
=== RUN   TestBroadcastActionsAndNoOps/stop-record
=== RUN   TestBroadcastActionsAndNoOps/pause-record
=== RUN   TestBroadcastActionsAndNoOps/resume-record
--- PASS: TestBroadcastActionsAndNoOps (0.01s)
    --- PASS: TestBroadcastActionsAndNoOps/start-stream (0.00s)
    --- PASS: TestBroadcastActionsAndNoOps/stop-stream (0.00s)
    --- PASS: TestBroadcastActionsAndNoOps/start-record (0.00s)
    --- PASS: TestBroadcastActionsAndNoOps/stop-record (0.00s)
    --- PASS: TestBroadcastActionsAndNoOps/pause-record (0.00s)
    --- PASS: TestBroadcastActionsAndNoOps/resume-record (0.00s)
=== RUN   TestBroadcastConcurrentDoubleSubmit
--- PASS: TestBroadcastConcurrentDoubleSubmit (0.00s)
=== RUN   TestBroadcastFreshPreflight
--- PASS: TestBroadcastFreshPreflight (0.00s)
=== RUN   TestBroadcastUnavailableRecording
=== RUN   TestBroadcastUnavailableRecording/false
=== RUN   TestBroadcastUnavailableRecording/true
--- PASS: TestBroadcastUnavailableRecording (0.00s)
    --- PASS: TestBroadcastUnavailableRecording/false (0.00s)
    --- PASS: TestBroadcastUnavailableRecording/true (0.00s)
=== RUN   TestBroadcastRejectedAndUnverified
=== RUN   TestBroadcastRejectedAndUnverified/rejection
=== RUN   TestBroadcastRejectedAndUnverified/pause_unsupported
=== RUN   TestBroadcastRejectedAndUnverified/ack_without_state
--- PASS: TestBroadcastRejectedAndUnverified (0.41s)
    --- PASS: TestBroadcastRejectedAndUnverified/rejection (0.00s)
    --- PASS: TestBroadcastRejectedAndUnverified/pause_unsupported (0.00s)
    --- PASS: TestBroadcastRejectedAndUnverified/ack_without_state (0.40s)
=== RUN   TestBroadcastDisconnectNeverReplays
--- PASS: TestBroadcastDisconnectNeverReplays (0.00s)
=== RUN   TestBroadcastValidationAndOrigin
--- PASS: TestBroadcastValidationAndOrigin (0.00s)
=== RUN   TestBroadcastOffline
--- PASS: TestBroadcastOffline (0.00s)
=== RUN   TestReadConfiguration
--- PASS: TestReadConfiguration (0.00s)
=== RUN   TestReadConfigurationRejectsInvalidSettings
--- PASS: TestReadConfigurationRejectsInvalidSettings (0.00s)
=== RUN   TestRunUsesConfiguredListeningAddress
2026/10/29 09:34:34 Controller: http://127.0.0.1:41933 — OBS bağlantısı tarayıcı durum isteğinde kurulur
--- PASS: TestRunUsesConfiguredListeningAddress (0.00s)
=== RUN   TestOBSAuthenticatedSnapshotAndScene
--- PASS: TestOBSAuthenticatedSnapshotAndScene (0.00s)
=== RUN   TestOfflineStatusAndSceneProtection
--- PASS: TestOfflineStatusAndSceneProtection (0.00s)
=== RUN   TestBitrateCounterResets
--- PASS: TestBitrateCounterResets (0.00s)
=== RUN   TestCancelledSnapshot
--- PASS: TestCancelledSnapshot (0.00s)
PASS
ok  	obs-lan-controller/Lib	0.537s
?   	obs-lan-controller/web	[no test files]
```
Limitations: race detector blocked by missing C compiler; real OBS mutations intentionally not verified.
