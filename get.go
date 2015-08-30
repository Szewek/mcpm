package main

import (
	"compress/bzip2"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type (
	AuthorInfo struct {
		Name, Url string
	}
	AttachmentInfo struct {
		Description, ThumbnailUrl, Title, Url string
		IsDefault                             bool
	}
	DependencyInfo struct {
		AddOnId, Type int
	}
	ModuleInfo struct {
		Foldername  string
		Fingerprint int64
	}
	FileInfo struct {
		Id                       int
		FileName, FileNameOnDisk string
		FileDate                 string
		ReleaseType, FileStatus  int
		DownloadURL              string
		IsAlternate              bool
		AlternateFileId          int
		Dependencies             []DependencyInfo
		IsAvailable              bool
		Modules                  []ModuleInfo
		PackageFingerprint       int64
		GameVersion              []string
	}
	CategoryInfo struct {
		Id        int
		Name, URL string
	}
	SectionInfo struct {
		ID, GameID              int
		Name                    string
		PackageType             int
		Path                    string
		InitialInclusionPattern string
		ExtraIncludePattern     string
	}
	GameFileInfo struct {
		GameVersion     string
		ProjectFileID   int
		ProjectFileName string
		FileType        int
	}
	ModInfo struct {
		Id                       int
		Name                     string
		Authors                  []AuthorInfo
		Attachments              []AttachmentInfo
		WebSiteURL               string
		GameId                   int
		Summary                  string
		DefaultFileId            int
		CommentCount             int
		DownloadCount            float64
		Rating                   int
		InstallCount             int
		IconId                   int
		LatestFiles              []FileInfo
		Categories               []CategoryInfo
		PrimaryAuthorName        string
		ExternalUrl              string
		Status                   int
		Stage                    int
		DonationUrl              string
		PrimaryCategoryId        int
		PrimaryCategoryName      string
		PrimaryCategoryAvatarUrl string
		Likes                    int
		CategorySection          SectionInfo
		PackageType              int
		AvatarUrl                string
		GameVersionLatestFiles   []GameFileInfo
		IsFeatured               int
		PopularityScore          float64
	}
	JSONResponse struct {
		Timestamp uint64    `json:"timestamp"`
		Data      []ModInfo `json:"data"`
	}
)

func updateCache() {
	// http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/{complete,hourly,daily,weekly}.json.bz2.txt
	resp, hte := http.Get("http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2")
	defer resp.Body.Close()
	if hte != nil {
		panic(hte)
	}
	fmt.Printf("BZIP size: %d\n", resp.ContentLength)
	bz := bzip2.NewReader(resp.Body)
	js := json.NewDecoder(bz)
	var jr JSONResponse
	jer := js.Decode(&jr)
	if jer != nil {
		panic(jer)
	}
	fmt.Printf("%v -- %v\n", len(jr.Data), jr.Data[0])
	gob.Register(AuthorInfo{})
	gob.Register(AttachmentInfo{})
	gob.Register(DependencyInfo{})
	gob.Register(ModuleInfo{})
	gob.Register(FileInfo{})
	gob.Register(CategoryInfo{})
	gob.Register(SectionInfo{})
	gob.Register(GameFileInfo{})
	gob.Register(ModInfo{})
	gob.Register(JSONResponse{})
	f, fe := os.OpenFile("file.gob.gz", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if fe != nil {
		panic(fe)
	}
	gz := gzip.NewWriter(f)
	defer gz.Close()
	gb := gob.NewEncoder(gz)
	ee := gb.Encode(&jr)
	if ee != nil {
		panic(ee)
	}
}
