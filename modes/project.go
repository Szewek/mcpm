package modes

import (
	"fmt"

	"github.com/Szewek/mcpm/mcpmdb"
)

func project(mo *ModeOptions) {
	pn := mo.Args[0]
	pkg := mcpmdb.GetPackage(pn)
	defer mcpmdb.Close()
	if pkg == nil {
		fmt.Printf("Project %#v doesn't exist!\n", pn)
		return
	}
	fmt.Printf("Project %#v information:\n", pn)
	pkg.PrintInfo()
}

func init() {
	registerMode("project", project)
}
