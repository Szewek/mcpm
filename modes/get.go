package modes

import (
	"fmt"
	"io"
	"os"

	"github.com/Szewek/mcpm/database"
	"github.com/Szewek/mcpm/util"
)

type (
	packageOptions struct {
		Dir          string
		ShouldUnpack bool
	}
)

var (
	pkgOptions = map[int]packageOptions{
		6: {"mods", false},
		5: {"", true},
		3: {"resourcepacks", false},
		1: {"saves", true},
	}
)

func get(mo *ModeOptions) {
	pn := mo.Args[0]
	db := database.GetDatabase()

	if en := db.Packages().GetByName(pn); en != nil {
		fn, pr, hte := util.DownloadPackage(en.Type, en.ID, en.Name, -1)

		util.Must(hte)
		defer pr.Close()

		fmt.Printf("Do you want to download package %#v? [y|N]", fn)
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

		dir := "."
		if pi, dok := pkgOptions[en.Type]; dok && pi.Dir != "" {
			dir = pi.Dir
		}
		util.Must(util.MkDirIfNotExist(dir))

		sav := fmt.Sprintf("%s/%s", dir, fn)
		fmt.Printf("Saving \"%s\"...\n", sav)
		f, fe := os.OpenFile(sav, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		util.Must(fe)

		_, ce := io.Copy(f, pr)
		util.Must(ce)
		fmt.Printf("Successfully saved to \"%s\"", sav)

		if pkgOptions[en.Type].ShouldUnpack {
			fmt.Println("This package should be unpacked in newer versions.")
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
