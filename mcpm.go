package main

import (
	"os"

	"github.com/Szewek/mcpm/modes"
)

var (
//	smodes = map[string]func(){
//		"get":        getPackage,
//		"search":     searchPackage,
//		"update":     updateCache,
//		"updatetest": func() { database.UpdateDatabase(true) },
//		"info":       readPackageInfo,
//	}
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
