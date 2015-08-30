package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	flagset     = flag.NewFlagSet("", flag.ExitOnError)
	forceUpdate bool
	modes       = map[string]func(){
		"update": updateCache,
	}
)

func main() {
	var mode string
	if len(os.Args) >= 2 {
		mode = os.Args[1]
	} else {
		mode = "help"
	}
	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, "mcpm – Minecraft Package Manager")
		fmt.Fprint(os.Stderr, "Available modes:\n ")
		for m, _ := range modes {
			fmt.Fprintf(os.Stderr, " %s", m)
		}
		fmt.Fprintln(os.Stderr, "\nAvailable options: ")
		flagset.PrintDefaults()
	}
	flagset.BoolVar(&forceUpdate, "f", false, "Forces update")
	if len(os.Args) >= 3 {
		flagset.Parse(os.Args[2:])
	}
	if f, ok := modes[mode]; ok {
		f()
	} else {
		flagset.Usage()
	}
}
