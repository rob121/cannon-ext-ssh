package main

import (
	"testing"

	"github.com/rob121/cannon/extension"
)

func TestCanUseTerminal(t *testing.T) {
	if canUseTerminal(nil) {
		t.Fatal("nil user")
	}
	if canUseTerminal(map[string]any{"authenticated": true}) {
		t.Fatal("missing permission")
	}
	user := map[string]any{
		"authenticated": true,
		"permissions":   []any{permFull},
	}
	if !canUseTerminal(user) {
		t.Fatal("expected allow")
	}
	req := extension.WireRequest{User: user}
	if !extension.UserCan(req, permFull) {
		t.Fatal("UserCan")
	}
}
