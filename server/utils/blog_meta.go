package utils

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ip2region "github.com/lionsoul2014/ip2region/binding/golang/service"
)

var (
	blogIPRegionOnce sync.Once
	blogIPRegion     *ip2region.Ip2Region
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
	if source := blogGeoIPSource(ip); source != "" {
		return source
	}
	return "public-network"
}

func blogGeoIPSource(ip string) string {
	searcher := blogIPRegionSearcher()
	if searcher == nil {
		return ""
	}
	region, err := searcher.Search(ip)
	if err != nil {
		return ""
	}
	return formatBlogIPRegion(region)
}

func blogIPRegionSearcher() *ip2region.Ip2Region {
	blogIPRegionOnce.Do(func() {
		v4Path := findBlogIPRegionDB("ip2region_v4.xdb")
		v6Path := findBlogIPRegionDB("ip2region_v6.xdb")
		if v4Path == "" && v6Path == "" {
			return
		}
		searcher, err := ip2region.NewIp2RegionWithPath(v4Path, v6Path)
		if err == nil {
			blogIPRegion = searcher
		}
	})
	return blogIPRegion
}

func findBlogIPRegionDB(name string) string {
	candidates := []string{
		filepath.Join("resource", "ip2region", name),
		filepath.Join("..", "resource", "ip2region", name),
		filepath.Join("server", "resource", "ip2region", name),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "resource", "ip2region", name))
	}
	for _, candidate := range candidates {
		if stat, err := os.Stat(candidate); err == nil && !stat.IsDir() {
			return candidate
		}
	}
	return ""
}

func formatBlogIPRegion(region string) string {
	parts := strings.Split(region, "|")
	formatted := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "0" {
			continue
		}
		formatted = append(formatted, part)
	}
	return strings.Join(formatted, " ")
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
