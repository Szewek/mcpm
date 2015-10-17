package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PackageOptions struct {
	Dir          string
	ShouldUnpack bool
}

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
)

func DownloadPackage(typ int, pid int, name string, fid int) (string, io.ReadCloser, error) {
	us := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/", pkgURLDirs[typ], pid, name)
	var us2 string
	if us2 = "latest"; fid != -1 {
		us2 = fmt.Sprintf("%d/download", fid)
	}
	download := fmt.Sprint(us, us2)
	ht, hte := http.Get(download)
	if hte != nil {
		return "", nil, hte
	}
	fname := ht.Request.URL.Path
	fname = fname[strings.LastIndex(fname, "/")+1:]
	return fname, NewProgressReader(ht.Body, uint64(ht.ContentLength), fmt.Sprintf("Downloading %#v (package %#v)...", fname, name)), nil
}

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
	return NewProgressReader(ht.Body, uint64(ht.ContentLength), fmt.Sprintf("Downloading %#v information...", fln)), nil
}

func GetPackageOptions(typ int) *PackageOptions {
	if k, g := pkgOptions[typ]; g {
		return &k
	}
	return &PackageOptions{".", false}
}
