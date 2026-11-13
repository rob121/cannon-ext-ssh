package main

import (
	"embed"
	"log"

	"github.com/rob121/cannon/extension"
)

//go:embed help/*.md
var helpFiles embed.FS

//go:embed static/xterm.js static/xterm.css
var staticFiles embed.FS

const (
	extensionName        = "cannon-ext-ssh"
	extensionTitle       = "SSH Terminal"
	extensionDescription = "Web terminal in Cannon admin (local shell as the extension process user)."
	extensionVersion     = "0.1.3"
	updateURLBase        = "https://github.com/rob121/cannon-ext-ssh/releases/download"
	defaultAdminMenu     = "SSH Terminal"
	permUse              = "terminal.use"
)

func main() {
	if err := newServer().Run(); err != nil {
		log.Fatal(err)
	}
}

func newServer() *extension.Server {
	s := extension.New(extension.Info{
		Name:          extensionName,
		Title:         extensionTitle,
		Description:   extensionDescription,
		Version:       extensionVersion,
		UpdateURLBase: updateURLBase,
		AdminMenuName: defaultAdminMenu,
	})
	s.EmbedHelp(helpFiles, "help")
	s.RegisterPermissions([]extension.PermissionDef{
		{
			ID:          permUse,
			DisplayName: "Use web terminal",
			Description: "Open an admin shell on the server.",
			Dangerous:   true,
		},
	})
	s.HandleAdmin("/admin", handleAdmin)
	s.Handle("/terminal/ws", handleTerminalWS)
	return s
}
