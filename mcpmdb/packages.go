package mcpmdb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Szewek/mcpm/helper"
	"github.com/Szewek/mcpm/util"
)

type (
	// MCPMFileList is a list of package files
	MCPMFileList struct {
		a []*MCPMFile
	}
	// MCPMFile represents a file which can be downloaded.
	MCPMFile struct {
		ID, Name, Release, MCVersion string
		Pkg                          *MCPMPackage
	}
)

// PrintInfo prints information about package.
func (pkg *MCPMPackage) PrintInfo() {
	fmt.Printf(" ID: %d\n Name: %s (%s)\n Type: %s\n", pkg.id, pkg.title, pkg.name, pkg.ptype)
}

// GetFileList creates list of files available for download.
func (pkg *MCPMPackage) GetFileList() *MCPMFileList {
	cur := fmt.Sprintf(util.CurseMCURL, strconv.FormatUint(pkg.id, 10))
	cht, chte := http.Get(cur)
	defer util.MustClose(cht.Body)
	defer util.Must(chte)
	cdoc, cdoce := goquery.NewDocumentFromResponse(cht)
	util.Must(cdoce)
	list := cdoc.Find("table.project-file-listing tbody tr")
	if list.Length() == 0 {
		return nil
	}
	fl := &MCPMFileList{make([]*MCPMFile, list.Length())}
	list.Each(func(i int, n *goquery.Selection) {
		td := n.Find("td")
		na := td.Eq(0).Find("a")
		nt := td.Eq(1)
		nv := td.Eq(2)
		tid := na.AttrOr("href", "")
		fl.a[i] = &MCPMFile{tid[strings.LastIndexByte(tid, '/')+1:], na.Text(), nt.Text(), nv.Text(), pkg}
	})
	return fl
}

// GetLatest returns first file entry from the list.
func (fl *MCPMFileList) GetLatest() *MCPMFile {
	return fl.a[0]
}

// Download saves the file to a file system and optionally unpacks its contents.
func (f *MCPMFile) Download() {
	pkgo := util.GetPackageOptionsB(f.Pkg.ptype)
	util.Must(util.MkDirIfNotExist(pkgo.Dir))
	sav := fmt.Sprintf("%s/%s", pkgo.Dir, f.Name)
	sf, sfe := os.OpenFile(sav, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	util.Must(sfe)
	defer util.MustClose(sf)
	us := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/%s/download", f.Pkg.ptype, f.Pkg.id, f.Pkg.name, f.ID)
	ht, hte := http.Get(us)
	util.Must(hte)
	pr := util.NewProgressReader(ht.Body, uint64(ht.ContentLength), fmt.Sprintf("Downloading %#v (package %#v)...", f.Name, f.Pkg.name))
	defer util.MustClose(pr)
	_, ce := io.Copy(sf, pr)
	util.Must(ce)
	fmt.Printf("Successfully saved to \"%s\"\n", sav)
	if pkgo.ShouldUnpack {
		switch f.Pkg.ptype {
		case "modpacks":
			helper.NewModPackHelper(sav).Unpack()
			fmt.Printf("Successfully installed modpack %#v\n", f.Pkg.title)
			break
		case "worlds":
			svh := helper.NewSaveHelper(sav)
			svh.UnpackAll()
			fmt.Printf("Successfully installed world save %#v\n", f.Pkg.title)
			break
		}
	}
}
