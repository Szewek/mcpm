package modes

import (
	"flag"
	"fmt"
	"os"
)

type (
	// ModeOptions contains all flags and arguments set in a command.
	ModeOptions struct {
		ModeName     string
		Verbose      bool
		DownloadOnly bool
		VersionQuery string
		Args         []string
	}
	// Mode contains functions needed in specific order. (EXPERIMENTAL)
	Mode struct {
		Func      func(*Mode, *ModeOptions)
		Usage     string
		RunBefore *Mode
		RunAfter  *Mode
	}
	// ModeList contains modes associated with unique name.
	ModeList map[string]func(*ModeOptions)
)

func (m *Mode) CanRun() bool {
	return m.Func != nil
}
func (m *Mode) Run(mo *ModeOptions) {
	if mb := m.RunBefore; mb != nil {
		mb.Run(mo)
	}
	m.Func(m, mo)
	if ma := m.RunAfter; ma != nil {
		ma.Run(mo)
	}
}

var (
	modelist  = &ModeList{}
	modelist2 = map[string]*Mode{}
	flg       = flag.NewFlagSet("", flag.ExitOnError)
)

// LaunchMode finds mode by its name and run this with set mode options.
// If mode is not found, LaunchMode shows usage of command.
func LaunchMode(m string) {
	mo := &ModeOptions{}
	mo.ModeName = m
	flg.Usage = func() {
		fmt.Fprintln(os.Stderr, "mcpm – Minecraft Package Manager\nAvailable modes:")
		for m := range *modelist {
			fmt.Fprintf(os.Stderr, "  %s\n", m)
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
	if f, ok := (*modelist)[m]; ok {
		f(mo)
	} else {
		flg.Usage()
	}
}

// EXPERIMENTAL
func runModeByName(m string) {
	if md, g := modelist2[m]; g {
		mo := &ModeOptions{}
		mo.ModeName = m
		md.Run(mo)
	}
}

func registerMode(m string, fn func(*ModeOptions)) {
	(*modelist)[m] = fn
}
