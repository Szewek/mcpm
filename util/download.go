package util

import (
	"fmt"
	"io"
	"net/http"
)

type (
	PackageOptions struct {
		Dir          string
		ShouldUnpack bool
	}
)

var (
	pkgURLDirs = map[int]string{
		6: "mc-mods",
		5: "modpacks",
		3: "texture-packs",
		1: "worlds",
	}
	pkgOptions = map[int]PackageOptions{
		6: {"mods", false},
		5: {".", true},
		3: {"resourcepacks", false},
		1: {"saves", true},
	}
	pkgTypeOptions = map[string]PackageOptions{
		"mc-mods":       {"mods", false},
		"modpacks":      {".", true},
		"texture-packs": {"resourcepacks", false},
		"worlds":        {"saves", true},
	}
)

//DownloadPackageInfo allows to download information about packages.
func DownloadPackageInfo(typ int, pid int, name string) (io.ReadCloser, error) {
	dp := fmt.Sprintf("http://widget.mcf.li/%s/minecraft/", pkgURLDirs[typ])
	var fln string
	if typ != 6 {
		fln = fmt.Sprintf("%d-%s", pid, name)
	} else {
		fln = name
	}
	dp = fmt.Sprint(dp, fln, ".json")
	ht, hte := http.Get(dp)
	if hte != nil {
		return nil, hte
	}
	fmt.Printf("Downloading %#v information...\n", fln)
	return NewReadProgress(ht.Body, uint64(ht.ContentLength)), nil
}

func GetPackageOptions(typ int) *PackageOptions {
	if k, g := pkgOptions[typ]; g {
		return &k
	}
	return &PackageOptions{".", false}
}
func GetPackageOptionsB(typ string) *PackageOptions {
	if k, g := pkgTypeOptions[typ]; g {
		return &k
	}
	return &PackageOptions{".", false}
}
