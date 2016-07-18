// MCPM is a command-line tool for managing Minecraft resources (packages)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

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

var (
	cpuprof = flag.String("cpuprof", "", "CPU Profile")
	memprof = flag.String("memprof", "", "Memory Profile")
)

func main() {
	flag.Parse()
	fmt.Print(intro)
	if *cpuprof != "" {
		fmt.Println("CPU profile is ON")
		f, err := os.Create(*cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	launchApp(flag.Args())
	if *memprof != "" {
		fmt.Println("Memory profile is ON")
		f, err := os.Create(*memprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

func launchApp(args []string) {
	var m string
	if len(args) > 0 {
		m = args[0]
	} else {
		m = "help"
	}
	modes.LaunchMode(m, args[1:])
}

func measureMemStats() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Println("Total allocs: ", ms.Mallocs)
	fmt.Println("Total frees: ", ms.Frees)
}
