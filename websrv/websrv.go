package websrv

import (
	"fmt"
	"net/http"
)

// LaunchWebServer starts a web server and listens on port 8080
func LaunchWebServer() {
	addr := ":8080"
	http.HandleFunc("/", handleHTML)
	http.HandleFunc("/js", handleJS)
	http.HandleFunc("/css", handleCSS)
	fmt.Printf("Listening on %#v...", addr)
	http.ListenAndServe(addr, nil)
}
