package database

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Szewek/mcpm/util"
)

const remoteURL = "http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2"

type (
	authorInfo struct {
		Name, Url string
	}
	categoryInfo struct {
		Id        int
		Name, URL string
	}
	modInfo struct {
		Id              int
		Name            string
		Authors         []authorInfo
		WebSiteURL      string
		Summary         string
		DownloadCount   float64
		Rating          int
		InstallCount    int
		Categories      []categoryInfo
		ExternalUrl     string
		Status          int
		Stage           int
		DonationUrl     string
		Likes           int
		PackageType     int
		IsFeatured      int
		PopularityScore float64
	}
	response struct {
		Timestamp uint64    `json:"timestamp"`
		Data      []modInfo `json:"data"`
	}
)

func UpdateDatabase(verbose bool) {
	fmt.Println("Updating database...")
	resp, hte := http.Get(remoteURL)
	defer util.Must(resp.Body.Close())
	util.Must(hte)
	if verbose {
		fmt.Printf("Bzipped JSON file size: %s", util.FileSize(resp.ContentLength))
	}

	pr := util.NewProgressReader(resp.Body, resp.ContentLength)
	bz := bzip2.NewReader(pr)
	js := json.NewDecoder(bz)

	var r response
	util.Must(js.Decode(r))

	db = &database{&pkglist{map[string]int{}, map[int]PkgElement{}}, &filelist{}}

	i := 0
	for ; i < len(r.Data); i++ {
		pkg := r.Data[i]
		indx := strings.LastIndexAny(pkg.WebSiteURL, "/") + 1
		var name string
		if indx > 0 {
			name = pkg.WebSiteURL[indx:]
		} else {
			name = fmt.Sprintf("--%d", pkg.Id)
			if verbose {
				fmt.Printf("Found unknown package: %#v; Naming it %#v", pkg.WebSiteURL, name)
			}
		}
		sid := fmt.Sprint(pkg.Id)
		idix := strings.Index(name, sid)
		if idix == 0 {
			name = name[len(sid)+1:]
		}
		authors := []string{}
		for j := 0; j < len(pkg.Authors); j++ {
			authors = append(authors, pkg.Authors[j].Name)
		}
		db.p.names[name] = pkg.Id
		db.p.pkgs[pkg.Id] = PkgElement{pkg.Id, pkg.PackageType, name, pkg.Name, pkg.Summary, authors}
	}
	fmt.Printf("Saving %d packages...\n", i)
	util.Must(util.WriteGobGzip(dbFile, db))
	fmt.Println("Database updated.")
}
