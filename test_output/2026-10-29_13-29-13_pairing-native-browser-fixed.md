# Native Chromium form acceptance
UTC: 2026-10-29T13:29:13.227541+00:00
Real running server, machine LAN IP; real OBS read-only, no mutations.
Before fix: native requestSubmit on /pair with dummy code produced Forbidden Origin. Cause: Referrer-Policy no-referrer causes native form Origin null, rejected by strict server.
Fix: pairing page same-origin referrer policy; strict Origin checks preserved, null still forbidden.
After rebuild/restart: local code read internally without printing; LAN native form submitted; navigation / and /api/deck status 200 connected true. Local native revoke form navigated / without Forbidden Origin; subsequent LAN API returned pairing_required (401). First immediate fetch raced navigation and failed; retry after navigation passed. No codes/cookies recorded. All previous sessions invalidated by server restart. Physical phone untested.
