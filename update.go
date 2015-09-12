package main

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type (
	_AuthorInfo struct {
		Name, Url string
	}
	_CategoryInfo struct {
		Id        int
		Name, URL string
	}
	_ModInfo struct {
		Id              int
		Name            string
		Authors         []_AuthorInfo
		WebSiteURL      string
		Summary         string
		DownloadCount   float64
		Rating          int
		InstallCount    int
		Categories      []_CategoryInfo
		ExternalUrl     string
		Status          int
		Stage           int
		DonationUrl     string
		Likes           int
		PackageType     int
		IsFeatured      int
		PopularityScore float64
	}
	_Response struct {
		Timestamp uint64     `json:"timestamp"`
		Data      []_ModInfo `json:"data"`
	}
)

func updateCache() {
	// http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/{complete,hourly,daily,weekly}.json.bz2.txt
	fmt.Println("Downloading database...")
	resp, hte := http.Get("http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2")
	defer resp.Body.Close()
	must(hte)
	if verbose {
		fmt.Printf("Bzipped JSON size: %s\n", fileSize(resp.ContentLength))
	}
	pr := newProgressReader(resp.Body, resp.ContentLength)
	bz := bzip2.NewReader(pr)
	js := json.NewDecoder(bz)
	var jr _Response
	jer := js.Decode(&jr)
	must(jer)
	db := &_Database{}
	pn := &_PackageNames{}
	i := 0
	for ; i < len(jr.Data); i++ {
		mod := jr.Data[i]
		indx := strings.LastIndexAny(mod.WebSiteURL, "/") + 1
		var pkgn string
		if indx > 0 {
			pkgn = mod.WebSiteURL[indx:]
		} else {
			pkgn = fmt.Sprintf("--%d", mod.Id)
			if verbose {
				fmt.Printf("Found unknown package: %#v; Naming it %#v", mod.WebSiteURL, pkgn)
			}
		}
		sid := fmt.Sprint(mod.Id)
		idix := strings.Index(pkgn, sid)
		if idix == 0 {
			pkgn = pkgn[len(sid)+1:]
		}
		authors := []string{}
		for j := 0; j < len(mod.Authors); j++ {
			authors = append(authors, mod.Authors[j].Name)
		}
		(*pn)[pkgn] = mod.Id
		(*db)[mod.Id] = _DataElement{mod.Id, _PackageType(mod.PackageType), pkgn, mod.Name, mod.Summary, authors}
	}
	fmt.Printf("Packages count: %d\nSaving...\n", i)
	writeGobGzip(homePath(dbFile), db)
	writeGobGzip(homePath(pnFile), pn)
	writeGob(homePath(luFile), time.Now())
	_DBASE = db
	_PKGS = pn
	fmt.Println("Database updated.")
}

func updatePackageCache() {}

// TODO Updating Forge versions' cache
func updateForgeCache() {}
