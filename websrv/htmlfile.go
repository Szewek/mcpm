package websrv

import (
	"fmt"
	"net/http"
	"strings"
)

var html = `<!DOCTYPE html><html><head>
<title>MCPM</title>
<link rel="stylesheet" href="/css" />
</head><body>
<header><menu-drawer-button></menu-drawer-button><h1>MCPM - Minecraft Package Manager</h1></header>
<main>Wait for future releases to see more configuration via web server...</main>
<script src="/js"></script>
</body></html>`

func handleHTML(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
}

func init() {
	rplc := strings.NewReplacer("\n", "")
	html = rplc.Replace(html)
}
