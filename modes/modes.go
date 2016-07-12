package modes

import (
	"flag"
	"fmt"
	"os"
)

type (
	// ModeOptions contains all flags and arguments set in a command.
	ModeOptions struct {
		ModeName string
		Verbose  bool
		Args     []string
	}
	ModeCommand interface {
		CanRun() bool
		Run(*ModeOptions)
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
	if m.CanRun() {
		m.Func(m, mo)
	}
	if ma := m.RunAfter; ma != nil {
		ma.Run(mo)
	}
}

var (
	modelist = &ModeList{}
	//modelist2 = map[string]*Mode{}
	flg = flag.NewFlagSet("", flag.ExitOnError)
)

// LaunchMode finds mode by its name and run this with set mode options.
// If mode is not found, LaunchMode shows usage of command.
func LaunchMode(m string) {
	mo := &ModeOptions{}
	mo.ModeName = m
	flg.Usage = func() {
		fmt.Fprintln(os.Stderr, "MCPM -- Minecraft Package Manager\nAvailable modes:")
		for m := range *modelist {
			fmt.Fprintf(os.Stderr, "  %s\n", m)
		}
		fmt.Fprintln(os.Stderr, "\nAvailable options: ")
		flg.PrintDefaults()
	}
	flg.BoolVar(&(mo.Verbose), "v", false, "Verbose (WIP)")
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
/*func runModeByName(m string) {
	if md, g := modelist2[m]; g {
		mo := &ModeOptions{}
		mo.ModeName = m
		md.Run(mo)
	}
}*/

func registerMode(m string, fn func(*ModeOptions)) {
	(*modelist)[m] = fn
}
