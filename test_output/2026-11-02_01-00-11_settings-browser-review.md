# Settings browser targeted test

Date: 2026-11-02T01:00:12+03:00

Purpose: browser interactions and production middleware composition.

Command: `go test ./Lib -run TestSettingsProductionChainBrowser -count=1` (Chrome and Node paths supplied).

Result: exit 0. Chromium against test server, simulated OBS; no real OBS connection. Physical devices and real OBS unverified.
