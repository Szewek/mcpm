package websrv

import (
	"fmt"
	"net/http"
	"strings"
)

var css = strings.TrimSpace(`
* {
	box-sizing: border-box;
}
body {
	margin: 0;
	background: #DFDFDF;
	font-family: sans-serif;
}
header {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	height: 48px;
	background: #505050;
	color: #FFF;
}
menu-drawer-button {
	display: inline-block;
	width: 32px;
	height: 32px;
	margin: 8px;
	border-radius: 50%;
	background: rgba(255,255,255,0.2);
}
h1 {
	display: inline-block;
	margin: 0;
	line-height: 48px;
	vertical-align: top;
}
main {
	margin-top: 48px;
	padding: 8px;
}
`)

func handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(w, css)
}

func init() {
	rplc := strings.NewReplacer("\n\t", "", "\n}", "}", ": ", ":")
	css = rplc.Replace(css)
}
