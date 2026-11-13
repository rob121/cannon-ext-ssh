package main

import (
	"strings"

	"github.com/rob121/cannon/extension"
)

// adminWebSocketURL builds an absolute WebSocket URL (tests/helpers).
// The in-browser terminal uses parent.location + path instead of this value.
func adminWebSocketURL(req extension.WireRequest) string {
	path := "/admin/extension-apps/" + extensionName + "/ws"
	abs := extension.PublicRequestURL(req, path)
	if strings.HasPrefix(abs, "https://") {
		return "wss://" + strings.TrimPrefix(abs, "https://")
	}
	if strings.HasPrefix(abs, "http://") {
		return "ws://" + strings.TrimPrefix(abs, "http://")
	}
	return path
}
