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

var (
	pkgURLDirs = map[_PackageType]string{
		type_Mod:          "mc-mods",
		type_ModPack:      "modpacks",
		type_ResourcePack: "texture-packs",
		type_WorldSave:    "worlds",
	}
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
		download := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/latest", pkgURLDirs[data.Type], pid, data.PkgName)
		if verbose {
			fmt.Printf("Checking URL %#v\n", download)
		}
		ht, hte := http.Get(download)
		must(hte)
		defer ht.Body.Close()
		if verbose {
			fmt.Printf("Found file URL: %#v\n", ht.Request.URL.String())
		}
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
		mkDirIfNotExist(dir)
		save := fmt.Sprintf("%s/%s", dir, fname)
		fmt.Printf("Saving \"%s\"...\n", save)
		f, fe := os.OpenFile(save, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		must(fe)
		lr := newProgressReader(ht.Body, ht.ContentLength)
		_, ce := io.Copy(f, lr)
		must(ce)
		fmt.Printf("Successfully saved to \"%s\"", save)
	} else {
		if verbose {
			fmt.Printf("Package %#v not found.\n", pkgn)
		}
		fmt.Println("What is that package?")
	}
}
