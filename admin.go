package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rob121/cannon/extension"
)

func handleAdmin(req extension.WireRequest) extension.WireResponse {
	if !extension.UserAuthenticated(req) || !extension.UserCan(req, permFull) {
		return extension.Error(http.StatusForbidden, "forbidden")
	}
	body, err := adminPageHTML()
	if err != nil {
		return extension.Error(http.StatusInternalServerError, "admin ui")
	}
	return extension.HTML(http.StatusOK, body)
}

func adminPageHTML() (string, error) {
	css, err := staticFiles.ReadFile("static/xterm.css")
	if err != nil {
		return "", err
	}
	js, err := staticFiles.ReadFile("static/xterm.js")
	if err != nil {
		return "", err
	}
	wsPath := "/admin/extension-apps/" + extensionName + "/ws"
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width,initial-scale=1\">")
	b.WriteString("<title>SSH Terminal</title><style>")
	b.Write(css)
	b.WriteString(`html,body{margin:0;height:100%;width:100%;overflow:hidden;background:#0d1117}
#terminal{position:absolute;inset:0;padding:4px;box-sizing:border-box}
#terminal .xterm,#terminal .xterm-viewport,#terminal .xterm-screen{height:100%!important;width:100%!important}
</style></head><body>`)
	b.WriteString("<div id=\"terminal\" aria-label=\"Shell terminal\"></div><script>")
	b.Write(js)
	b.WriteString(fmt.Sprintf(`</script><script>
(function(){
  var term = new Terminal({ cursorBlink: true, theme: { background: "#0d1117", foreground: "#e6edf3" } });
  var el = document.getElementById("terminal");
  term.open(el);
  function fit(){
    var cols = Math.max(20, Math.floor(el.clientWidth / 9));
    var rows = Math.max(5, Math.floor(el.clientHeight / 17));
    term.resize(cols, rows);
    return { cols: cols, rows: rows };
  }
  var sz = fit();
  var wsPath = %q;
  function resolveWS(){
    try {
      var p = parent.location;
      return (p.protocol === "https:" ? "wss:" : "ws:") + "//" + p.host + wsPath;
    } catch (e) {}
    return (location.protocol === "https:" ? "wss:" : "ws:") + "//" + location.host + wsPath;
  }
  var ws = new WebSocket(resolveWS());
  ws.binaryType = "arraybuffer";
  ws.onopen = function(){
    ws.send(JSON.stringify({ type: "resize", cols: sz.cols, rows: sz.rows }));
    term.focus();
  };
  ws.onerror = function(){
    term.write("\r\n\r\n[connection failed — is Cannon running and is the extension started?]\r\n");
  };
  ws.onmessage = function(ev){
    if (ev.data instanceof ArrayBuffer) {
      term.write(new Uint8Array(ev.data));
    } else {
      term.write(ev.data);
    }
  };
  ws.onclose = function(ev){
    if (ev.code !== 1000) {
      term.write("\r\n\r\n[session closed: " + ev.code + "]\r\n");
    } else {
      term.write("\r\n\r\n[session closed]\r\n");
    }
  };
  term.onData(function(data){ if (ws.readyState === 1) ws.send(data); });
  window.addEventListener("resize", function(){
    sz = fit();
    if (ws.readyState === 1) {
      ws.send(JSON.stringify({ type: "resize", cols: sz.cols, rows: sz.rows }));
    }
  });
})();
</script></body></html>`, wsPath))
	return b.String(), nil
}
