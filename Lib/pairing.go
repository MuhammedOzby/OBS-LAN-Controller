package lib

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
	"obs-lan-controller/web"
)

const (
	pairingCodeTTL       = 5 * time.Minute
	pairingSessionTTL    = 12 * time.Hour
	pairingMaxSessions   = 128
	pairingMaxFailures   = 20
	pairingFailureWindow = time.Minute
	pairingCookieName    = "obs_controller_session"
)

// State is deliberately process-local: restarting invalidates every code/session.
// One pending code and a global failure window avoid attacker-keyed allocations.
type pairingHandler struct {
	next         http.Handler
	mu           sync.Mutex
	now          func() time.Time
	code         string
	codeExpires  time.Time
	sessions     map[string]time.Time
	failures     int
	failureReset time.Time
	listenAddr   net.Addr
}

type pairingContextKey struct{}
type pairingQR struct {
	Address string
	Image   template.URL
}
type pairingView struct {
	Code    string
	Expires string
	Error   string
	QRs     []pairingQR
	Local   bool
}

var pairBootstrap, _ = web.Assets.ReadFile("static/pair-bootstrap.js")

// Destinations derive only from the bound listener and local interfaces, never HTTP headers.
func pairingAddresses(addr net.Addr) []string {
	if addr == nil {
		return nil
	}
	tcp, ok := addr.(*net.TCPAddr)
	if !ok || tcp.Port == 0 {
		return nil
	}
	bound := tcp.IP
	if bound != nil && bound.IsLoopback() {
		return nil
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var out []string
	for _, a := range addrs {
		ip, _, err := net.ParseCIDR(a.String())
		if err != nil || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() || ip.IsMulticast() {
			continue
		}
		if bound != nil && !bound.IsUnspecified() && !bound.Equal(ip) {
			continue
		}
		if bound != nil && bound.To4() != nil && ip.To4() == nil {
			continue
		}
		out = append(out, net.JoinHostPort(ip.String(), strconv.Itoa(tcp.Port)))
	}
	sort.Strings(out)
	return out
}

func (p *pairingHandler) view(local, rotate bool) (pairingView, error) {
	v := pairingView{Local: local}
	if !local {
		return v, nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	p.prune(now)
	if rotate || p.code == "" {
		code, err := pairingRandom(6)
		if err != nil {
			return v, err
		}
		p.code, p.codeExpires = code, now.Add(pairingCodeTTL)
	}
	v.Code = p.code
	v.Expires = p.codeExpires.Local().Format("15:04:05")
	for _, address := range pairingAddresses(p.listenAddr) {
		png, err := qrcode.Encode("http://"+address+"/pair#"+v.Code, qrcode.Medium, 256)
		if err != nil {
			return pairingView{}, err
		}
		v.QRs = append(v.QRs, pairingQR{address, template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(png))})
	}
	return v, nil
}

func newPairingHandler(next http.Handler) *pairingHandler {
	return &pairingHandler{next: next, now: time.Now, sessions: make(map[string]time.Time)}
}

func (p *pairingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pair-bootstrap.js" {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", 405)
			return
		}
		w.Write(pairBootstrap)
		return
	}
	if r.URL.Path == "/pair" {
		p.pair(w, r)
		return
	}
	if pairingLoopback(r) || p.authenticated(r) {
		if r.URL.Path == "/mobil-cihaz-ile-giris" && r.Method == http.MethodGet {
			if (len(r.Header.Values("Origin")) != 0 && !pairingSameOrigin(r)) || r.Header.Get("Sec-Fetch-Site") == "cross-site" {
				http.Error(w, "Forbidden Origin", 403)
				return
			}
			v, err := p.view(pairingLoopback(r), false)
			if err != nil {
				http.Error(w, "Pairing unavailable", 503)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), pairingContextKey{}, v))
		}
		p.next.ServeHTTP(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"pairing_required"}`))
		return
	}
	http.Redirect(w, r, "/pair", http.StatusSeeOther)
}

// Trust the socket peer only, never Host, proxy headers, or DNS names.
func pairingLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(host)
	return err == nil && ip != nil && ip.IsLoopback()
}

func pairingSameOrigin(r *http.Request) bool {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	origins := r.Header.Values("Origin")
	if len(origins) != 1 {
		return false
	}
	origin, err := url.Parse(origins[0])
	return err == nil && origin.Scheme == scheme && origin.Host == r.Host && origin.User == nil && origin.Path == "" && origin.RawQuery == "" && !origin.ForceQuery && origin.Fragment == "" && origin.Opaque == ""
}

func (p *pairingHandler) prune(now time.Time) {
	for token, expiry := range p.sessions {
		if !now.Before(expiry) {
			delete(p.sessions, token)
		}
	}
	if !now.Before(p.codeExpires) {
		p.code = ""
	}
	if !now.Before(p.failureReset) {
		p.failures = 0
		p.failureReset = now.Add(pairingFailureWindow)
	}
}

func (p *pairingHandler) authenticated(r *http.Request) bool {
	cookie, err := r.Cookie(pairingCookieName)
	if err != nil || len(cookie.Value) != 64 {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	p.prune(now)
	expiry, ok := p.sessions[cookie.Value]
	return ok && now.Before(expiry)
}

func pairingRandom(bytes int) (string, error) {
	raw := make([]byte, bytes)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(raw)), nil
}

func (p *pairingHandler) pair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	// Native form POSTs under no-referrer can carry Origin: null. Preserve
	// same-origin provenance while still omitting cross-origin referrers.
	w.Header().Set("Referrer-Policy", "same-origin")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; img-src data:; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Remote GET only renders the public form and can be entered from a scanner
	// application. Local GET issues a secret and every POST changes state, so
	// both retain the strict Origin and cross-site checks.
	protected := r.Method == http.MethodPost || pairingLoopback(r)
	if protected && (((r.Method == http.MethodPost || len(r.Header.Values("Origin")) != 0) && !pairingSameOrigin(r)) || r.Header.Get("Sec-Fetch-Site") == "cross-site") {
		http.Error(w, "Forbidden Origin", http.StatusForbidden)
		return
	}
	var submitted string
	if r.Method == http.MethodPost {
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid pairing request", http.StatusBadRequest)
			return
		}
		if len(r.PostForm["code"]) == 1 {
			submitted = r.PostForm.Get("code")
		}
	}
	if r.Method == http.MethodGet {
		v, err := p.view(pairingLoopback(r), true)
		if err != nil {
			http.Error(w, "Pairing unavailable", 503)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		pairingPage.Execute(w, v)
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	p.prune(now)
	if r.Method == http.MethodPost && r.PostForm.Get("action") == "revoke-all" {
		if !pairingLoopback(r) {
			http.Error(w, "Local access required", http.StatusForbidden)
			return
		}
		clear(p.sessions)
		p.code = ""
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if p.failures >= pairingMaxFailures {
		w.Header().Set("Retry-After", "60")
		http.Error(w, "Too many pairing attempts; try again later", http.StatusTooManyRequests)
		return
	}
	if p.code == "" || subtle.ConstantTimeCompare([]byte(submitted), []byte(p.code)) != 1 {
		p.failures++
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		pairingPage.Execute(w, pairingView{Error: "Kod geçersiz, kullanılmış veya süresi dolmuş. OBS bilgisayarında yeni QR kodu açın veya yeni kodu elle girin."})
		return
	}
	if len(p.sessions) >= pairingMaxSessions {
		http.Error(w, "Session limit reached; wait for expiry or restart locally", http.StatusServiceUnavailable)
		return
	}
	token, err := pairingRandom(32)
	if err != nil {
		http.Error(w, "Pairing unavailable", 503)
		return
	}
	expiry := now.Add(pairingSessionTTL)
	p.sessions[token] = expiry
	p.code = "" // Consume only after session creation succeeds.
	http.SetCookie(w, &http.Cookie{Name: pairingCookieName, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil, MaxAge: int(pairingSessionTTL / time.Second), Expires: expiry})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var pairingPage = template.Must(template.Must(template.New("pair.tmpl").Parse(`<!doctype html>
<html lang="tr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Controller · Eşleştirme</title>
<style>:root{color-scheme:dark;font-family:system-ui,sans-serif;background:#111315;color:#f5f0e5}*{box-sizing:border-box}body{margin:0;min-height:100dvh;display:grid;place-items:center;padding:1rem}main{width:min(100%,30rem);min-width:0;background:#1b1e22;border:1px solid #454039;border-radius:1rem;padding:clamp(1rem,4vw,2rem)}h1{font-size:1.5rem;color:#ffc107}p{line-height:1.6;color:#d1cec7;overflow-wrap:anywhere}.pair-qr-list img{max-width:100%;height:auto}output{display:block;font-size:clamp(1.3rem,6vw,2rem);letter-spacing:.15em;font-weight:700;color:#ffc107;overflow-wrap:anywhere}label{display:block;margin-bottom:.5rem}input,button{font:inherit;width:100%;border-radius:.5rem;padding:.8rem}input{background:#111315;color:#fff;border:1px solid #76683c}button{margin-top:1rem;background:#ffc107;color:#17140a;border:0;font-weight:700;cursor:pointer}a{color:#ffc107}:focus-visible{outline:3px solid #ffc107;outline-offset:3px}</style></head>
<body><main><h1>Controller eşleştirme</h1>{{template "pair-content" .}}<p><a href="/">Panele dön</a></p></main><script src="/pair-bootstrap.js" defer></script></body></html>`)).ParseFS(web.Assets, "templates/pair-content.tmpl"))
