package main

import (
	"strings"
	"testing"

	"github.com/rob121/cannon/extension"
)

func TestAdminWebSocketURL(t *testing.T) {
	req := extension.WireRequest{
		URL: "/admin/extension-apps/cannon-ext-ssh",
		Settings: map[string]string{
			extension.SettingSitePublicBase: "http://127.0.0.1:8001",
		},
	}
	got := adminWebSocketURL(req)
	want := "ws://127.0.0.1:8001/admin/extension-apps/cannon-ext-ssh/ws"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	req.Settings[extension.SettingSitePublicBase] = "https://example.com"
	got = adminWebSocketURL(req)
	if !strings.HasPrefix(got, "wss://example.com/") {
		t.Fatalf("https base: %q", got)
	}
}
