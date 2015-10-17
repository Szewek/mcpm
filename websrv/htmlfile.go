package websrv

import (
	"fmt"
	"net/http"
)

var html = `<!DOCTYPE html><html><head>
<title>MCPM</title>
</head><body>
<header><h1>MCPM</h1></header>
<main></main>
</body></html>`

func handleHTML(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html)
}
