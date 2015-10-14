package main

import (
	"os"

	"github.com/Szewek/mcpm/modes"
)

func main() {
	var m string
	if len(os.Args) >= 2 {
		m = os.Args[1]
	} else {
		m = "help"
	}
	modes.LaunchMode(m)
}
