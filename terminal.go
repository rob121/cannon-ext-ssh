package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/rob121/cannon/extension"
)

const permFull = extensionName + "." + permUse

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type resizeMsg struct {
	Type string `json:"type"`
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

func canUseTerminal(user map[string]any) bool {
	if auth, ok := user["authenticated"].(bool); !ok || !auth {
		return false
	}
	return extension.UserCan(extension.WireRequest{User: user}, permFull)
}

func handleTerminalWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, ok := extension.UserFromProxyHeader(r.Header)
	if !ok || !canUseTerminal(user) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	cmd, ptmx, err := startShellPTY()
	if err != nil {
		log.Printf("pty start: %v", err)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(shellErrMessage(err)))
		return
	}
	defer func() {
		_ = ptmx.Close()
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	var closed atomic.Bool
	go func() {
		buf := make([]byte, 4096)
		for !closed.Load() {
			n, err := ptmx.Read(buf)
			if n > 0 {
				if werr := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); werr != nil {
					break
				}
			}
			if err != nil {
				break
			}
		}
		closed.Store(true)
	}()

	for !closed.Load() {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		switch msgType {
		case websocket.TextMessage:
			if len(data) > 0 && data[0] == '{' {
				var m resizeMsg
				if json.Unmarshal(data, &m) == nil && m.Type == "resize" && m.Cols > 0 && m.Rows > 0 {
					_ = pty.Setsize(ptmx, &pty.Winsize{Cols: m.Cols, Rows: m.Rows})
					continue
				}
			}
			if _, werr := ptmx.Write(data); werr != nil {
				break
			}
		case websocket.BinaryMessage:
			if _, werr := ptmx.Write(data); werr != nil {
				break
			}
		}
	}
	closed.Store(true)
	_, _ = io.Copy(io.Discard, ptmx)
}
