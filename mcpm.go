// MCPM is a command-line tool for managing Minecraft resources (packages)
package main

import (
	"os"

	"github.com/Szewek/mcpm/modes"
	"github.com/Szewek/mcpm/websrv"
)

func main() {
	var m string
	if len(os.Args) >= 2 {
		m = os.Args[1]
	} else {
		m = "help"
	}

	if m == "websrv" {
		websrv.LaunchWebServer()
	} else {
		modes.LaunchMode(m)
	}
}
