package modes

import (
	"flag"
	"fmt"
	"os"
)

type (
	ModeOptions struct {
		Verbose      bool
		DownloadOnly bool
		VersionQuery string
		Args         []string
	}
	ModeList map[string]func(*ModeOptions)
)

var (
	modelist = ModeList{}
	flg      = flag.NewFlagSet("", flag.ExitOnError)
)

func LaunchMode(m string) {
	mo := &ModeOptions{}
	flg.Usage = func() {
		fmt.Fprintln(os.Stderr, "mcpm – Minecraft Package Manager\nAvailable modes:\n ")
		for m, _ := range modelist {
			fmt.Fprintf(os.Stderr, " %s\n", m)
		}
		fmt.Fprintln(os.Stderr, "\nAvailable options: ")
		flg.PrintDefaults()
	}
	flg.BoolVar(&(mo.Verbose), "v", false, "Verbose (WIP)")
	flg.BoolVar(&(mo.DownloadOnly), "d", false, "Download only")
	flg.StringVar(&(mo.VersionQuery), "q", "", "Version query (like \"latest:beta:mc1.7.10\")")
	if len(os.Args) >= 3 {
		flg.Parse(os.Args[2:])
		mo.Args = flg.Args()
	}
	if f, ok := modelist[m]; ok {
		f(mo)
	} else {
		flg.Usage()
	}
}
