package websrv

import (
	"fmt"
	"net/http"
	"strings"
)

var js = `
`

func handleJS(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, js)
}

func init() {
	rplc := strings.NewReplacer("\n", "")
	js = rplc.Replace(js)
}
