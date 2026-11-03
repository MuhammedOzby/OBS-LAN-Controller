# Actual running server pairing acceptance
UTC: 2026-10-29T13:14:03.258105+00:00
Real TCP via machine LAN address (not RemoteAddr fixture); same-machine client, no physical phone. OBS read-only. Code/cookie never logged.
```json
{
  "unpaired_api": 401,
  "unpaired_page_redirect": true,
  "pairing_cookie_count": 1,
  "paired_deck_readback": true,
  "program": "Konuşma",
  "single_use": 401,
  "revoked_readback": 401,
  "final_outputs": {
    "currentScene": "Konuşma",
    "streamActive": false,
    "recordActive": false
  }
}
```
All assertions passed. Test session revoked. HTTP unencrypted; TLS/PWA remains separate.