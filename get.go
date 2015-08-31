package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	dbFile = ".mcpmdb"
	pnFile = ".mcpmpn"
	luFile = ".mcpmlu"
)

// TODO Get version from http://minecraft.curseforge.com/mc-mods/<id>-<pkgname>/files/<fileid>/download
func getPackage() {
	pkgn := flagset.Arg(0)
	pn, pne := readPackageNamesFromFile(pnFile)
	must(pne)
	db, dbe := readDatabaseFromFile(dbFile)
	must(dbe)
	if pid, ok := (*pn)[pkgn]; ok {
		data := db.Data[pid]
		download := fmt.Sprintf("http://minecraft.curseforge.com/mc-mods/%d-%s/files/latest", pid, data.PkgName)
		ht, hte := http.Get(download)
		must(hte)
		defer ht.Body.Close()
		fname := ht.Request.URL.Path
		fname = fname[strings.LastIndex(fname, "/")+1:]
		fmt.Printf("Do you want to download package %#v? [y|N]", fname)
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
		if _, der := os.Stat(dir); os.IsNotExist(der) {
			must(os.Mkdir(dir, 438))
		}
		save := fmt.Sprintf("%s/%s", dir, fname)
		fmt.Printf("Saving \"%s\"...\n", save)
		f, fe := os.OpenFile(save, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		must(fe)
		_, ce := io.Copy(f, ht.Body)
		must(ce)
		fmt.Printf("Successfully saved to \"%s\"", save)
	} else {
		fmt.Println("What is that package?")
	}
}
