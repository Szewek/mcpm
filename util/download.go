package util

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	pkgURLDirs = map[int]string{
		6: "mc-mods",
		5: "modpacks",
		3: "texture-packs",
		1: "worlds",
	}
)

func DownloadPackage(typ int, pid int, name string, fid int) (string, *ProgressReader, error) {
	us := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/", pkgURLDirs[typ], pid, name)
	var us2 string
	if us2 = "latest"; fid != -1 {
		us2 = fmt.Sprintf("%d/download", fid)
	}
	download := fmt.Sprint(us, us2)
	fmt.Printf("Checking URL %#v\n", download)
	ht, hte := http.Get(download)
	if hte != nil {
		return "", nil, hte
	}
	fmt.Printf("Found file URL: %#v\n", ht.Request.URL.String())
	fname := ht.Request.URL.Path
	fname = fname[strings.LastIndex(fname, "/")+1:]
	return fname, NewProgressReader(ht.Body, ht.ContentLength), nil
}
func DownloadPackageInfo(typ int, pid int, name string) (*ProgressReader, error) {
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
	return NewProgressReader(ht.Body, ht.ContentLength), nil
}
