package main

import (
	"fmt"
	"io"
	"os"
)

const (
	luFile = ".mcpmlu"
	pnFile = ".mcpmpn"
	dbFile = ".mcpmdb"
)

type (
	_PackageOptions struct {
		Dir          string
		ShouldUnpack bool
	}
)

var (
	pkgOptions = map[_PackageType]_PackageOptions{
		type_Mod:          {"mods", false},
		type_ModPack:      {"", true},
		type_ResourcePack: {"resourcepacks", false},
		type_WorldSave:    {"saves", true},
	}
)

// TODO Get version from http://minecraft.curseforge.com/<type>/<id>-<pkgname>/files/<fileid>/download
func getPackage() {
	pkgn := flagset.Arg(0)
	if pid, ok := (*_PKGS)[pkgn]; ok {
		data := (*_DBASE)[pid]
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
