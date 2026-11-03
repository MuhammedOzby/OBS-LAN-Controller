package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lib "obs-lan-controller/Lib"

	"github.com/gin-gonic/gin"
)

func TestHandlerUsesGin(t *testing.T) {
	handler, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := handler.(*gin.Engine); !ok {
		t.Fatalf("handler type = %T, want *gin.Engine", handler)
	}
}

func TestDesignPages(t *testing.T) {
	handler, err := lib.NewHandler()
	if err != nil {
		t.Fatal(err)
	}
	for _, page := range []struct{ path, title string }{
		{"/", "Kontrol"},
		{"/deck", "Deck"},
		{"/yayin-yonetimi", "Yayın Yönetimi"},
	} {
		t.Run(page.path, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, page.path, nil))
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d", rr.Code)
			}
			body := rr.Body.String()
			if strings.Contains(body, `class="app-header"`) || strings.Contains(body, `role="banner"`) {
				t.Error("global application header should be removed; section/editor headers are allowed")
			}
			navStart, navEnd := strings.Index(body, "<nav"), strings.Index(body, "</nav>")
			if navStart < 0 || navEnd < navStart || !strings.Contains(body[navStart:navEnd], "connection-status") {
				t.Error("connection status should be inside navigation")
			}
			if strings.Contains(body, `class="page-heading"`) || strings.Contains(body, `class="design-label"`) {
				t.Error("page heading and design label should be removed")
			}
			for _, want := range []string{"<title>" + page.title + " · OBS LAN Controller</title>", "OBS bağlı değil", "/static/vendor/bootstrap.min.css", "aria-current=\"page\""} {
				if !strings.Contains(body, want) {
					t.Errorf("missing %q", want)
				}
			}
			if got := strings.Count(body, "class=\"nav-link"); got != 4 {
				t.Errorf("navigation items = %d, want 4", got)
			}
		})
	}
}
