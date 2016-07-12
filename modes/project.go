package modes

import (
	"fmt"

	"github.com/Szewek/mcpm/util"
)

func project(mo *ModeOptions) {
	pn := mo.Args[0]
	xpi := util.GetCurseProjectInfo(pn)
	if xpi == nil {
		fmt.Printf("Project %#v doesn't exist!\n", pn)
		return
	}
	fmt.Printf("Project %#v information:\n", pn)
	fmt.Printf(" ID: %s\n Name: %s\n Type: %s\n Newest file: %s\n", xpi.ID, xpi.Name, xpi.Type, xpi.Files[0].Name)
}

func init() {
	registerMode("project", project)
}
