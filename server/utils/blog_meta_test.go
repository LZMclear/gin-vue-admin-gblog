package utils

import "testing"

func TestBlogIPSourceSpecialNetworks(t *testing.T) {
	if got := BlogIPSource("127.0.0.1"); got != "localhost" {
		t.Fatalf("BlogIPSource loopback = %q, want localhost", got)
	}
	if got := BlogIPSource("192.168.1.10"); got != "private-network" {
		t.Fatalf("BlogIPSource private = %q, want private-network", got)
	}
}

func TestBlogIPSourcePublicGeoIP(t *testing.T) {
	got := BlogIPSource("219.133.110.197")
	if got == "" || got == "public-network" || got == "unknown" {
		t.Fatalf("BlogIPSource public = %q, want geo location", got)
	}
}
