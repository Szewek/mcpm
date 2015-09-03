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
	_AttachmentInfo struct {
		Description, ThumbnailUrl, Title, Url string
		IsDefault                             bool
	}
	_DependencyInfo struct {
		AddOnId, Type int
	}
	_ModuleInfo struct {
		Foldername  string
		Fingerprint int64
	}
	_FileInfo struct {
		Id                       int
		FileName, FileNameOnDisk string
		FileDate                 string
		ReleaseType, FileStatus  int
		DownloadURL              string
		IsAlternate              bool
		AlternateFileId          int
		Dependencies             []_DependencyInfo
		IsAvailable              bool
		Modules                  []_ModuleInfo
		PackageFingerprint       int64
		GameVersion              []string
	}
	_CategoryInfo struct {
		Id        int
		Name, URL string
	}
	_SectionInfo struct {
		ID, GameID              int
		Name                    string
		PackageType             int
		Path                    string
		InitialInclusionPattern string
		ExtraIncludePattern     string
	}
	_GameFileInfo struct {
		GameVersion     string
		ProjectFileID   int
		ProjectFileName string
		FileType        int
	}
	_ModInfo struct {
		Id                       int
		Name                     string
		Authors                  []_AuthorInfo
		Attachments              []_AttachmentInfo
		WebSiteURL               string
		GameId                   int
		Summary                  string
		DefaultFileId            int
		CommentCount             int
		DownloadCount            float64
		Rating                   int
		InstallCount             int
		IconId                   int
		LatestFiles              []_FileInfo
		Categories               []_CategoryInfo
		PrimaryAuthorName        string
		ExternalUrl              string
		Status                   int
		Stage                    int
		DonationUrl              string
		PrimaryCategoryId        int
		PrimaryCategoryName      string
		PrimaryCategoryAvatarUrl string
		Likes                    int
		CategorySection          _SectionInfo
		PackageType              int
		AvatarUrl                string
		GameVersionLatestFiles   []_GameFileInfo
		IsFeatured               int
		PopularityScore          float64
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
	bz := bzip2.NewReader(resp.Body)
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
		versions := []_FileElement{}
		for j := 0; j < len(mod.Authors); j++ {
			authors = append(authors, mod.Authors[j].Name)
		}
		for j := 0; j < len(mod.GameVersionLatestFiles); j++ {
			fi := mod.GameVersionLatestFiles[j]
			versions = append(versions, _FileElement{fi.ProjectFileID, fi.GameVersion})
		}
		(*pn)[pkgn] = mod.Id
		(*db)[mod.Id] = _DataElement{mod.Id, PackageType(mod.PackageType), pkgn, mod.Name, mod.Summary, authors, versions}
	}
	fmt.Printf("Packages count: %d\nSaving...\n", i)
	db.Save(homePath(dbFile))
	pn.Save(homePath(pnFile))
	writeGob(homePath(luFile), time.Now())
	fmt.Println("Database updated.")
}
