package lib

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
)

// httpRequestReadTimeout bounds headers plus body reads for every route, including
// action JSON decoding. Change this server policy here, not the OBS configuration.
const httpRequestReadTimeout = 5 * time.Second

// newHTTPServer is the production boundary: router-only fixtures intentionally
// remain usable with synthetic hosts. No forwarded header is a trust source.
func newHTTPServer(addr string, handler http.Handler) (*http.Server, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("HTTP trusted hosts: enumerate local interfaces: %w", err)
	}
	local := make(map[netip.Addr]bool, len(addrs))
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, fmt.Errorf("HTTP trusted hosts: invalid interface address")
		}
		if ip, ok := netip.AddrFromSlice(ip); ok && !ip.IsUnspecified() {
			local[ip.Unmap()] = true
		}
	}
	guard := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !trustedHTTPHost(r.Host, local) {
			http.Error(w, "Untrusted Host", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	})
	return &http.Server{
		Addr: addr, Handler: guard,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       httpRequestReadTimeout,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}, nil
}

// trustedHTTPHost does not resolve names: DNS resolving to this machine is not
// proof of trust. Interface IPs are an exact startup snapshot, not local subnets.
func trustedHTTPHost(authority string, local map[netip.Addr]bool) bool {
	host := authority
	if strings.HasPrefix(authority, "[") && strings.HasSuffix(authority, "]") {
		host = authority[1 : len(authority)-1]
		ip, err := netip.ParseAddr(host)
		if err != nil || !ip.Is6() {
			return false
		}
	} else if strings.Contains(authority, ":") {
		var port string
		var err error
		host, port, err = net.SplitHostPort(authority)
		if err != nil || port == "" || strings.IndexFunc(port, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
			return false
		}
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return false
		}
	}
	if strings.EqualFold(host, "localhost") || strings.EqualFold(host, "localhost.") {
		return true
	}
	ip, err := netip.ParseAddr(host)
	if err != nil || ip.Zone() != "" {
		return false
	}
	ip = ip.Unmap()
	return ip.IsLoopback() || (!ip.IsUnspecified() && local[ip])
}
