package main

import (
	"fmt"
	"io"
	"os"
)

const (
	dbFile = ".mcpmdb"
	pnFile = ".mcpmpn"
	luFile = ".mcpmlu"
)

// TODO Get version from http://minecraft.curseforge.com/<type>/<id>-<pkgname>/files/<fileid>/download
func getPackage() {
	pkgn := flagset.Arg(0)
	pn, pne := readPackageNamesFromFile(homePath(pnFile))
	must(pne)
	db, dbe := readDatabaseFromFile(homePath(dbFile))
	must(dbe)
	if pid, ok := (*pn)[pkgn]; ok {
		data := (*db)[pid]
		fn, pr, hte := downloadPackage(&data, -1)
		must(hte)
		defer pr.Close()
		fmt.Printf("Do you want to download package %#v? [y|N]", fn)
		r := []byte{}
		_, se := fmt.Scanln(&r)
		must(se)
		if len(r) == 0 {
			fmt.Println("Download cancelled by default")
			return
		}
		if r[0] != 0x79 && r[0] != 0x59 {
			fmt.Println("Cancelled file download.")
			return
		}
		var dir string
		switch data.Type {
		case type_WorldSave:
			dir = "saves"
			break
		case type_ResourcePack:
			dir = "resourcepacks"
			break
		case type_Mod:
			dir = "mods"
			break
		default:
			dir = "."
		}
		mkDirIfNotExist(dir)
		save := fmt.Sprintf("%s/%s", dir, fn)
		fmt.Printf("Saving \"%s\"...\n", save)
		f, fe := os.OpenFile(save, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		must(fe)
		_, ce := io.Copy(f, pr)
		must(ce)
		fmt.Printf("Successfully saved to \"%s\"", save)
	} else {
		if verbose {
			fmt.Printf("Package %#v not found.\n", pkgn)
		}
		fmt.Println("What is that package?")
	}
}

// TODO Downloading Forge
func getForge() {}
