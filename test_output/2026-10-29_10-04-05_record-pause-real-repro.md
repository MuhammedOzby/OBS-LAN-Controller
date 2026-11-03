# Gerçek OBS kayıt repro
UTC: 2026-10-29T10:04:05.724770+00:00
Kullanıcı izniyle HTTP API start-record, pause-record, resume-record; ardından GET durum okuma. Yayın değiştirilmedi.
```json
[
  {
    "initial": {
      "connected": true,
      "recordAvailable": true,
      "recordActive": false,
      "recordPaused": false,
      "recordDuration": 0,
      "recordBytes": 37292262,
      "streamActive": false
    }
  },
  {
    "action": "start-record",
    "http": 200,
    "state": {
      "connected": true,
      "recordAvailable": true,
      "recordActive": true,
      "recordPaused": false,
      "recordDuration": 533,
      "recordBytes": 37758220,
      "streamActive": false
    }
  },
  {
    "action": "pause-record",
    "http": 502,
    "state": {
      "connected": true,
      "recordAvailable": true,
      "recordActive": true,
      "recordPaused": false,
      "recordDuration": 1566,
      "recordBytes": 38497133,
      "streamActive": false
    }
  },
  {
    "action": "resume-record",
    "http": 409,
    "state": {
      "connected": true,
      "recordAvailable": true,
      "recordActive": true,
      "recordPaused": false,
      "recordDuration": 2566,
      "recordBytes": 39267613,
      "streamActive": false
    }
  }
]
```
