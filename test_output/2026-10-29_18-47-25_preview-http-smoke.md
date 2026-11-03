# Preview production HTTP smoke
UTC 2026-10-29_18-47-25
Environment: isolated XDG_RUNTIME_DIR, example empty-password fixture only. No real OBS.

HTTP /api/preview/status idle: {"connected": false, "viewers": 0, "streaming": false, "width": 0, "height": 0, "fps": 0, "error": ""}
Served embedded HTML/JS/CSS match changed source
Built production binary: same-UID Unix handshake + HTTP readback passed (simulated plugin)
SIGTERM clean shutdown and socket cleanup passed
