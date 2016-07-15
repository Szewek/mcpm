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
	}
)

// PrintInfo prints information about package.
func (pkg *MCPMPackage) PrintInfo() {
	fmt.Printf(" ID: %d\n Name: %s (%s)\n Type: %s\n", pkg.id, pkg.title, pkg.name, pkg.ptype)
}

// GetFileList creates list of files available for download.
func (pkg *MCPMPackage) GetFileList() *MCPMFileList {
	cur := fmt.Sprintf("http://curse.com/project/%s", strconv.FormatUint(pkg.id, 10))
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
		fl.a[i] = &MCPMFile{tid[strings.LastIndexByte(tid, '/')+1:], na.Text(), nt.Text(), nv.Text()}
	})
	return fl
}

// DownloadFileWithID downloads file with a specified ID and unpacks its contents if necessary.
func (pkg *MCPMPackage) DownloadFileWithID(fid string, buf []byte) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	po := util.GetPackageOptionsB(pkg.ptype)
	util.Must(util.MkDirIfNotExist(po.Dir))
	us := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/%s/download", pkg.ptype, pkg.id, pkg.name, fid)
	ht, hte := http.Get(us)
	util.Must(hte)
	fname := ht.Request.URL.Path
	fname = fname[strings.LastIndex(fname, "/")+1:]
	sav := fmt.Sprintf("%s/%s", po.Dir, fname)
	sf, sfe := os.OpenFile(sav, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	util.Must(sfe)
	defer util.MustClose(sf)
	pr := util.NewProgressReader(ht.Body, uint64(ht.ContentLength), fmt.Sprintf("Downloading \"%s\": %s", pkg.name, fname))
	defer util.MustClose(pr)
	_, ce := io.CopyBuffer(sf, pr, buf)
	util.Must(ce)
	fmt.Printf("Successfully saved to \"%s\"\n", sav)
	if po.ShouldUnpack {
		switch pkg.ptype {
		case "modpacks":
			NewModPackHelper(sav).Unpack()
			fmt.Printf("Successfully installed modpack %#v\n", pkg.title)
			break
		case "worlds":
			svh := helper.NewSaveHelper(sav)
			svh.UnpackAll()
			fmt.Printf("Successfully unpacked world save %#v\n", pkg.title)
			break
		}
	}
}

// GetLatest returns first file entry from the list.
func (fl *MCPMFileList) GetLatest() *MCPMFile {
	return fl.a[0]
}

// GetFromID finds file with a specified ID.
func (fl *MCPMFileList) GetFromID(id string) *MCPMFile {
	for _, f := range fl.a {
		if f.ID == id {
			return f
		}
	}
	return nil
}
