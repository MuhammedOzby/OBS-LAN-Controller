package lib

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"obs-lan-controller/web"

	"github.com/gin-gonic/gin"
)

// NewHandler creates the Gin router with embedded templates and static files.
func NewHandler() (http.Handler, error) {
	return newHandler(nil)
}

func newHandler(obs *OBSClient) (http.Handler, error) {
	return newHandlerWithPreview(obs, nil)
}

func newHandlerWithPreview(obs *OBSClient, preview *previewService) (http.Handler, error) {
	store, _ := loadFavorites("")
	return newHandlerWithFavorites(obs, preview, store)
}

func newHandlerWithFavorites(obs *OBSClient, preview *previewService, store *favoriteStore) (http.Handler, error) {
	if obs != nil {
		obs.favorites = store
	}
	pages, err := template.ParseFS(web.Assets, "templates/*.tmpl")
	if err != nil {
		return nil, err
	}
	static, err := fs.Sub(web.Assets, "static")
	if err != nil {
		return nil, err
	}
	router := gin.New()
	router.Use(gin.Recovery())
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	router.StaticFS("/static", http.FS(static))
	if preview != nil {
		preview.routes(router)
	}
	router.GET("/api/status", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if obs == nil {
			c.JSON(200, OBSStatus{Scenes: []string{}})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		status := obs.Snapshot(ctx, c.Query("preview") == "1")
		config := store.get()
		collection := ""
		if status.Connected && config != nil {
			obs.mu.Lock()
			if obs.conn != nil {
				collection, _ = obs.collectionLocked(ctx)
			}
			obs.mu.Unlock()
		}
		c.JSON(200, struct {
			OBSStatus
			Favorites [3]string `json:"favorites"`
		}{status, favoriteSlots(config, collection, status.Scenes)})
	})
	favoriteRoutes(router, obs, store)
	router.GET("/api/deck", deckHandler(obs))
	router.POST("/api/deck/action", deckActionHandler(obs))
	audioRoutes(router, obs)
	router.POST("/api/broadcast", broadcastHandler(obs))
	router.POST("/api/scene", func(c *gin.Context) {
		origin, err := url.Parse(c.GetHeader("Origin"))
		expectedScheme := "http"
		if c.Request.TLS != nil {
			expectedScheme = "https"
		}
		if err != nil || origin.Scheme != expectedScheme || origin.Host != c.Request.Host || origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
			c.Status(http.StatusForbidden)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 4096)
		var request struct {
			SceneName string `json:"sceneName"`
		}
		if c.ShouldBindJSON(&request) != nil || request.SceneName == "" {
			c.Status(400)
			return
		}
		if obs == nil {
			c.Status(503)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
		defer cancel()
		if err := obs.SetScene(ctx, request.SceneName); err != nil {
			if errors.Is(err, errSceneNotAllowed) {
				c.Status(400)
				return
			}
			c.Status(502)
			return
		}
		c.Status(204)
	})
	for _, page := range []struct{ Path, Title, Template string }{
		{"/", "Kontrol", "control"},
		{"/deck", "Deck", "deck"},
		{"/yayin-yonetimi", "Yayın Yönetimi", "broadcast"},
		{"/mobil-cihaz-ile-giris", "Mobil cihaz ile giriş", "mobile"},
	} {
		router.GET(page.Path, func(c *gin.Context) {
			var body bytes.Buffer
			data := struct {
				Path, Title, Template string
				Local                 bool
				Pair                  pairingView
			}{page.Path, page.Title, page.Template, pairingLoopback(c.Request), pairingView{}}
			if v, ok := c.Request.Context().Value(pairingContextKey{}).(pairingView); ok {
				data.Pair = v
			}
			if err := pages.ExecuteTemplate(&body, "layout", data); err != nil {
				log.Printf("Şablon hatası: %v", err)
				c.String(http.StatusInternalServerError, "Sayfa oluşturulamadı")
				return
			}
			c.Header("Cache-Control", "no-store")
			c.Data(http.StatusOK, "text/html; charset=utf-8", body.Bytes())
		})
	}
	return router, nil
}

// Run parses CLI arguments and serves until interrupted or a server error occurs.
func Run(args []string) error {
	if len(args) > 0 && args[0] == "serve" {
		args = args[1:]
	}
	flags := flag.NewFlagSet("obs-lan-controller", flag.ContinueOnError)
	configPath := flags.String("config", "configuration.json", "Yapılandırma dosyasının yolu")
	addr := flags.String("addr", "", "HTTP dinleme adresini yapılandırma yerine geçici olarak kullan")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("bilinmeyen komut: %s", flags.Arg(0))
	}
	config, err := ReadConfiguration(*configPath)
	bootstrap := errors.Is(err, os.ErrNotExist)
	if err != nil && !bootstrap {
		return err
	}
	if bootstrap {
		config.OBS.IP, config.OBS.Port = "127.0.0.1", 4455
		config.Server.ListeningIP, config.Server.Port = "127.0.0.1", 8080
		if *addr != "" {
			host, _, e := net.SplitHostPort(*addr)
			if e != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
				return errors.New("ilk kurulum yalnızca loopback adresinde çalışabilir")
			}
		}
	}
	if *addr == "" {
		*addr = net.JoinHostPort(config.Server.ListeningIP, strconv.Itoa(config.Server.Port))
	}
	obs := NewOBSClient(config)
	defer obs.Close()
	preview := newPreviewService()
	defer preview.close()
	if err := preview.start(); err != nil {
		log.Printf("Canlı önizleme kullanılamıyor (OBS kontrolü devam eder): %v", err)
	} else {
		log.Printf("Önizleme: %s/obs-lan-controller/preview.sock", os.Getenv("XDG_RUNTIME_DIR"))
	}
	store, err := loadFavorites(filepath.Join(filepath.Dir(*configPath), filepath.Base(*configPath)+".favorites.json"))
	if err != nil {
		return fmt.Errorf("favori ayarları okunamadı: %w", err)
	}
	handler, err := newHandlerWithFavorites(obs, preview, store)
	if err != nil {
		return err
	}
	handler = newSettingsHandler(handler, obs, *configPath, bootstrap)
	pairing := newPairingHandler(handler)
	preview.authorize = func(r *http.Request) bool { return pairingLoopback(r) || pairing.authenticated(r) }
	server, err := newHTTPServer(*addr, pairing)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	done := make(chan error, 1)
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	pairing.listenAddr = listener.Addr()
	go func() { done <- server.Serve(listener) }()
	log.Printf("Controller: http://%s — OBS bağlantısı tarayıcı durum isteğinde kurulur", *addr)
	select {
	case err := <-done:
		if err != http.ErrServerClosed {
			return err
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
	return nil
}
