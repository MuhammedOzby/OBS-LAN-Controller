# frontend-structure-red

UTC: 2026-10-29T08:50:52.895392+00:00

Komut: /usr/local/go/bin/go test . -run TestFrontendControlStructure -v

Çıkış: 1

```text
=== RUN   TestFrontendControlStructure
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] GET    /                         --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /deck                     --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
[GIN-debug] GET    /yayin-yonetimi           --> obs-lan-controller/Lib.NewHandler.func1 (2 handlers)
    frontend_test.go:19: control missing "id=\"control-dashboard\""
    frontend_test.go:19: control missing "id=\"program-preview\""
    frontend_test.go:19: control missing "id=\"preview-placeholder\""
    frontend_test.go:19: control missing "Geçici WebSocket ekran görüntüsü"
    frontend_test.go:19: control missing "Canlı video değildir"
    frontend_test.go:19: control missing "id=\"stream-state\""
    frontend_test.go:19: control missing "Yayın yapılmıyor"
    frontend_test.go:19: control missing "id=\"obs-fps\""
    frontend_test.go:19: control missing "id=\"obs-memory\""
    frontend_test.go:19: control missing "id=\"obs-render\""
    frontend_test.go:19: control missing "id=\"obs-skips\""
    frontend_test.go:19: control missing "id=\"stat-cpu\""
    frontend_test.go:19: control missing "id=\"stat-network\""
    frontend_test.go:19: control missing "id=\"stat-drop\""
    frontend_test.go:19: control missing "id=\"stat-quality\""
    frontend_test.go:19: control missing "id=\"scene-feedback\""
    frontend_test.go:19: control missing "CPU"
    frontend_test.go:19: control missing "AĞ"
    frontend_test.go:19: control missing "DROP"
    frontend_test.go:19: control missing "Kalite"
    frontend_test.go:21: scene slots = 0, want exactly 3
--- FAIL: TestFrontendControlStructure (0.00s)
FAIL
FAIL	obs-lan-controller	0.005s
FAIL
```

Sınır: Gerçek OBS yok. Önceki ./... RED ayrıca eşzamanlı backend eksik sembollerinden etkilendi; ilk rapor kabuk backtick hatası nedeniyle çıktısız kaldı. Raporlama tekrarında Python quoting hatası oldu; test çalışmadı. Buradaki kök paket RED gerçek çıktı kanıtıdır.
