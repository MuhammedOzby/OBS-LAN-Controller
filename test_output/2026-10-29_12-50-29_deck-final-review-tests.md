# Deck son kontrol — otomatik doğrulama

Tarih: 2026-10-29T12:50:29.827207+00:00

Amaç: Deck ve mevcut uygulama regresyonları. Ortam: Go/Node; WebSocket fixture simülasyonu. Gerçek OBS üzerinde komut çalıştırılmaz.

## go test ./... -count=1
Çıkış kodu: 127
```text
/bin/sh: line 1: go: command not found

```

## go vet ./...
Çıkış kodu: 127
```text
/bin/sh: line 1: go: command not found

```

## node --check web/static/deck.js
Çıkış kodu: 0
```text

```

Sınırlamalar: Gerçek OBS Deck işlemleri, fiziksel mobil cihaz, eşleştirme ve HTTPS doğrulanmadı.