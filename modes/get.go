package modes

import (
	"fmt"
	"io"
	"os"

	"github.com/Szewek/mcpm/database"
	"github.com/Szewek/mcpm/helper"
	"github.com/Szewek/mcpm/util"
)

func get(mo *ModeOptions) {
	pn := mo.Args[0]
	db := database.GetDatabase()

	if en := db.Packages().GetByName(pn); en != nil {
		fn, pr, hte := util.DownloadPackage(en.Type, en.ID, en.Name, -1)

		util.Must(hte)
		defer util.MustClose(pr)

		fmt.Printf("Do you want to download %#v? [y|N] ", fn)
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

		pkgo := util.GetPackageOptions(en.Type)
		util.Must(util.MkDirIfNotExist(pkgo.Dir))

		sav := fmt.Sprintf("%s/%s", pkgo.Dir, fn)
		f, fe := os.OpenFile(sav, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		util.Must(fe)
		defer util.MustClose(f)

		_, ce := io.Copy(f, pr)
		util.Must(ce)
		fmt.Printf("Successfully saved to \"%s\"\n", sav)

		if pkgo.ShouldUnpack {
			fmt.Println("This package should be unpacked in newer versions.")
			if en.Type == 5 {
				helper.NewModPackHelper(sav).Unpack()
				fmt.Printf("Successfully installed modpack %#v\n", en.FullName)
			}
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
