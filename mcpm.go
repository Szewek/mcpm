// MCPM is a command-line tool for managing Minecraft resources (packages)
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Szewek/mcpm/modes"
)

const intro = `
    XX
  XXXXXX      MCPM
XXXXXXXXXX
++XXXXXX--    Minecraft Package Manager
++++XX----
+++++-----
  +++---
    +-
`

func main() {
	fmt.Print(intro)
	var m string
	if len(os.Args) >= 2 {
		m = os.Args[1]
	} else {
		m = "help"
	}
	modes.LaunchMode(m)
	measureMemStats()
}

func measureMemStats() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("Total allocs: ", ms.Mallocs)
	fmt.Println("Total frees: ", ms.Frees)
}
