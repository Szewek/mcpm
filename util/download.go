package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	PackageOptions struct {
		Dir          string
		ShouldUnpack bool
	}
	ProjectInfo struct {
		ID, Name, Type string
		Files          []FileInfo
	}
	FileInfo struct {
		ID, Name, Release, MCVersion string
	}
)

const (
	_CurseForgeURL = "http://minecraft.curseforge.com/projects/%s"
	_CurseMCURL    = "http://curse.com/project/%s"
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
)

// DownloadPackage allows to download packages.
// If File ID (fid) equals -1, then it downloads the latest version.
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
	return NewProgressReader(ht.Body, uint64(ht.ContentLength), fmt.Sprintf("Downloading %#v information...", fln)), nil
}

func GetPackageOptions(typ int) *PackageOptions {
	if k, g := pkgOptions[typ]; g {
		return &k
	}
	return &PackageOptions{".", false}
}

// GetCurseProjectInfo downloads information about specified project.
// Project's name comes from Curse Project URL.
func GetCurseProjectInfo(s string) *ProjectInfo {
	ur := fmt.Sprintf(_CurseForgeURL, s)
	ht, hte := http.Get(ur)
	defer MustClose(ht.Body)
	Must(hte)
	doc, doce := goquery.NewDocumentFromResponse(ht)
	Must(doce)
	tit := doc.Find("h1.project-title a .overflow-tip").Text()
	oa := doc.Find(".RootGameCategory a").Eq(0)
	ida := doc.Find(".view-on-curse a").Eq(0)
	cat := oa.AttrOr("href", " ")
	id := ida.AttrOr("href", "")
	idn := id[strings.LastIndexByte(id, '/')+1:]
	if idn == "" {
		return nil
	}
	cur := fmt.Sprintf(_CurseMCURL, idn)
	cht, chte := http.Get(cur)
	defer MustClose(cht.Body)
	defer Must(chte)
	cdoc, cdoce := goquery.NewDocumentFromResponse(cht)
	Must(cdoce)
	list := cdoc.Find("table.project-file-listing tbody tr")
	xfi := make([]FileInfo, list.Length())
	list.Each(func(i int, n *goquery.Selection) {
		td := n.Find("td")
		na := td.Eq(0).Find("a")
		nt := td.Eq(1)
		nv := td.Eq(2)
		tid := na.AttrOr("href", "")
		xfi[i] = FileInfo{tid[strings.LastIndexByte(tid, '/')+1:], na.Text(), nt.Text(), nv.Text()}
	})
	return &ProjectInfo{idn, tit, cat[1:], xfi}
}
