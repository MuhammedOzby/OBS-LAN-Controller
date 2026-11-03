# Gerçek OBS yayın durumu ve performans

Tarih: 2026-10-29T08:43:18.838058+00:00

Ortam: Linux, Python websockets. Hedef: 192.168.88.254:4455.

Adımlar: OBS v5 kimlik doğrulaması; salt okunur GetStreamStatus (iki örnek, yaklaşık 3 saniye), GetStats, GetRecordStatus, GetVideoSettings. Bit hızı gönderilen bayt farkından ölçüldü; sayaç sıfırlanması/pasif yayın durumunda hesaplanmaz.

```json
{
  "stream": {
    "outputActive": true,
    "outputBytes": 85035773,
    "outputCongestion": 0.0,
    "outputDuration": 110233,
    "outputReconnecting": false,
    "outputSkippedFrames": 0,
    "outputTimecode": "00:01:50.233",
    "outputTotalFrames": 3307,
    "outputSkippedFramesPercent": 0.0
  },
  "stats": {
    "activeFps": 30.000000300000007,
    "availableDiskSpace": 699420.05078125,
    "averageFrameRenderTime": 1.747693,
    "cpuUsage": 0.5538638269939098,
    "memoryUsage": 709.74609375,
    "outputSkippedFrames": 0,
    "outputTotalFrames": 3322,
    "renderSkippedFrames": 26,
    "renderTotalFrames": 19142,
    "webSocketSessionIncomingMessages": 4,
    "webSocketSessionOutgoingMessages": 4,
    "renderSkippedFramesPercent": 0.1358,
    "outputSkippedFramesPercent": 0.0
  },
  "record": {
    "outputActive": false,
    "outputBytes": 0,
    "outputDuration": 0,
    "outputPaused": false,
    "outputTimecode": "00:00:00.000"
  },
  "video": {
    "baseHeight": 600,
    "baseWidth": 1024,
    "fpsDenominator": 1,
    "fpsNumerator": 30,
    "outputHeight": 600,
    "outputWidth": 1024
  },
  "sample_seconds": 3.004,
  "stream_bitrate_mbps": 6.164
}
```

Hatalar: Yok

Çıkış kodu: 0

Gerçek OBS kullanıldı. OBS durumunda değişiklik yapılmadı. Kısa ölçüm uzun süreli yayın kararlılığını, izleyici deneyimini veya GPU kullanımını doğrulamaz; Go uygulamasına entegrasyon henüz yok.
