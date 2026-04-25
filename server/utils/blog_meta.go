package utils

import (
	"net"
	"net/http"
	"strings"
)

func BlogClientIP(r *http.Request) string {
	headers := []string{
		"X-Real-IP",
		"X-Forwarded-For",
		"Proxy-Client-IP",
		"WL-Proxy-Client-IP",
		"HTTP_CLIENT_IP",
		"HTTP_X_FORWARDED_FOR",
	}
	for _, key := range headers {
		v := strings.TrimSpace(r.Header.Get(key))
		if v == "" || strings.EqualFold(v, "unknown") {
			continue
		}
		parts := strings.Split(v, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
		return v
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func BlogIPSource(ip string) string {
	if ip == "" {
		return "unknown"
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return "unknown"
	}
	if parsed.IsLoopback() {
		return "localhost"
	}
	if isPrivateIP(parsed) {
		return "private-network"
	}
	return "public-network"
}

func BlogParseUserAgent(ua string) (os string, browser string) {
	l := strings.ToLower(ua)
	switch {
	case strings.Contains(l, "windows"):
		os = "Windows"
	case strings.Contains(l, "android"):
		os = "Android"
	case strings.Contains(l, "iphone"), strings.Contains(l, "ipad"), strings.Contains(l, "ios"):
		os = "iOS"
	case strings.Contains(l, "mac os"), strings.Contains(l, "macintosh"):
		os = "macOS"
	case strings.Contains(l, "linux"):
		os = "Linux"
	default:
		os = "Unknown"
	}

	switch {
	case strings.Contains(l, "edg/"), strings.Contains(l, "edge/"):
		browser = "Edge"
	case strings.Contains(l, "chrome/"):
		browser = "Chrome"
	case strings.Contains(l, "firefox/"):
		browser = "Firefox"
	case strings.Contains(l, "safari/") && !strings.Contains(l, "chrome/"):
		browser = "Safari"
	case strings.Contains(l, "micromessenger"):
		browser = "WeChat"
	default:
		browser = "Unknown"
	}
	return
}

func isPrivateIP(ip net.IP) bool {
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}
	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
