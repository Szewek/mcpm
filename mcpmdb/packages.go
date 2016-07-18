package mcpmdb

import (
	"fmt"
	"io"
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
	cht, chte := get(util.DirPathJoin("http://curse.com/project", strconv.FormatUint(pkg.id, 10)))
	defer util.MustClose(cht.Body)
	defer util.Must(chte)
	cdoc, cdoce := goquery.NewDocumentFromResponse(cht)
	util.Must(cdoce)
	list := cdoc.Find("table.project-file-listing tbody tr")
	if list.Length() == 0 {
		return nil
	}
	fl := &MCPMFileList{make([]*MCPMFile, list.Length())}
	for i := range list.Nodes {
		n := list.Eq(i)
		td := n.Find("td")
		na := td.Eq(0).Find("a")
		nt := td.Eq(1)
		nv := td.Eq(2)
		tid := na.AttrOr("href", "")
		fl.a[i] = &MCPMFile{tid[strings.LastIndexByte(tid, '/')+1:], na.Text(), nt.Text(), nv.Text()}
	}
	return fl
}

// DownloadFileWithID downloads file with a specified ID and unpacks its contents if necessary.
func (pkg *MCPMPackage) DownloadFileWithID(fid string, buf []byte) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	po := util.GetPackageOptionsB(pkg.ptype)
	util.Must(util.MkDirIfNotExist(po.Dir))
	us := fmt.Sprintf("http://minecraft.curseforge.com/projects/%s/files/%s/download", pkg.name, fid)
	ht, hte := get(us)
	defer util.MustClose(ht.Body)
	util.Must(hte)
	fname := ht.Request.URL.Path
	fname = fname[strings.LastIndex(fname, "/")+1:]
	sav := util.DirPathJoin(po.Dir, fname)
	sf, sfe := os.OpenFile(sav, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	util.Must(sfe)
	defer util.MustClose(sf)
	fmt.Printf("Downloading \"%s\": %s\n", pkg.name, fname)
	pr := util.NewReadProgress(ht.Body, uint64(ht.ContentLength))
	_, ce := io.CopyBuffer(sf, pr, buf)
	util.Must(ce)
	util.MustClose(pr)
	fmt.Printf("Successfully saved to \"%s\"\n", sav)
	if po.ShouldUnpack {
		switch pkg.ptype {
		case "modpacks":
			UnpackModpack(sav)
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
