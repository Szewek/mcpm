package modes

import (
	"fmt"

	"github.com/Szewek/mcpm/mcpmdb"
	"github.com/Szewek/mcpm/util"
)

func get(mo *ModeOptions) {
	pn := mo.Args[0]
	defer mcpmdb.Close()
	if pkg := mcpmdb.GetPackage(pn); pkg != nil {
		if fl := pkg.GetFileList(); fl != nil {
			pf := fl.GetLatest()
			if !mo.Confirm {
				fmt.Printf("Do you want to download %#v? [y|N] ", pf.Name)
				r := []byte{}
				_, se := fmt.Scanln(&r)
				util.Must(se)
				if len(r) == 0 {
					fmt.Println("Download cancelled by default")
					return
				}
				if r[0] != 0x79 && r[0] != 0x59 {
					fmt.Println("Cancelled file download.")
					return
				}
			}
			pkg.DownloadFileWithID(pf.ID, nil)
		} else {
			fmt.Println("Package exists but has no files to download.")
		}
	} else {
		if mo.Verbose {
			fmt.Printf("Package %#v not found.\n", pn)
		}
		fmt.Println("What is that package?")
	}
}

func init() {
	registerMode("get", get)
}
