package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

var (
	flagset = flag.NewFlagSet("", flag.ExitOnError)
	verbose bool
	modes   = map[string]func(){
		"get":    getPackage,
		"update": updateCache,
	}
	homeDir = "."
)

func checkHomeDir() {
	u, ue := user.Current()
	must(ue)
	homeDir = fmt.Sprintf("%s/.mcpm", u.HomeDir)
	must(mkDirIfNotExist(homeDir))
}

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
	flagset.BoolVar(&verbose, "v", false, "Verbose (WIP)")
	if len(os.Args) >= 3 {
		flagset.Parse(os.Args[2:])
	}
	if f, ok := modes[mode]; ok {
		checkHomeDir()
		f()
	} else {
		flagset.Usage()
	}
}
