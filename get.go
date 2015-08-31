package main

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	DBFile = ".mcpmdb"
	PNFile = ".mcpmpn"
	LUFile = ".mcpmlu"
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

// TODO Get version from http://minecraft.curseforge.com/mc-mods/<id>-<pkgname>/files/<fileid>/download
func getPackage() {
	pkgn := flagset.Arg(0)
	pn, pne := readPackageNamesFromFile(PNFile)
	Must(pne)
	db, dbe := readDatabaseFromFile(DBFile)
	Must(dbe)
	if pid, ok := (*pn)[pkgn]; ok {
		data := db.Data[pid]
		download := fmt.Sprintf("http://minecraft.curseforge.com/mc-mods/%d-%s/files/latest", pid, data.PkgName)
		ht, hte := http.Get(download)
		Must(hte)
		defer ht.Body.Close()
		fname := ht.Request.URL.Path
		fname = fname[strings.LastIndex(fname, "/")+1:]
		fmt.Printf("Do you want to download package %#v? [y|N]", fname)
		r := []byte{}
		_, se := fmt.Scanln(&r)
		Must(se)
		if len(r) == 0 {
			fmt.Println("Download cancelled by default")
			return
		}
		if r[0] != 0x79 && r[0] != 0x59 {
			fmt.Println("Cancelled file download.")
			return
		}
		var dir string
		switch data.Type {
		case TypeWorldSave:
			dir = "saves"
			break
		case TypeResourcePack:
			dir = "resourcepacks"
			break
		case TypeMod:
			dir = "mods"
			break
		default:
			dir = "."
		}
		if _, der := os.Stat(dir); os.IsNotExist(der) {
			Must(os.Mkdir(dir, 438))
		}
		save := fmt.Sprintf("%s/%s", dir, fname)
		fmt.Printf("Saving \"%s\"...\n", save)
		f, fe := os.OpenFile(save, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
		Must(fe)
		_, ce := io.Copy(f, ht.Body)
		Must(ce)
		fmt.Printf("Successfully saved to \"%s\"", save)
	} else {
		fmt.Println("What is that package?")
	}
}
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
	db := createDatabase()
	pn := &PackageNames{}
	i := 0
	for ; i < len(jr.Data); i++ {
		mod := jr.Data[i]
		indx := strings.LastIndexAny(mod.WebSiteURL, "/") + 1
		var pkgn string
		if indx > 0 {
			pkgn = mod.WebSiteURL[indx:]
		} else {
			pkgn = fmt.Sprintf("--%d", mod.Id)
		}
		sid := fmt.Sprint(mod.Id)
		idix := strings.Index(pkgn, sid)
		if idix == 0 {
			pkgn = pkgn[len(sid)+1:]
		}
		authors := []string{}
		versions := []FileElement{}
		for j := 0; j < len(mod.Authors); j++ {
			authors = append(authors, mod.Authors[j].Name)
		}
		for j := 0; j < len(mod.GameVersionLatestFiles); j++ {
			fi := mod.GameVersionLatestFiles[j]
			versions = append(versions, FileElement{fi.ProjectFileID, fi.GameVersion})
		}
		(*pn)[pkgn] = mod.Id
		db.Data[mod.Id] = DataElement{mod.Id, PackageType(mod.PackageType), pkgn, mod.Name, mod.Summary, authors, versions}
	}
	fmt.Printf("Packages added: %d\nSaving...\n", i)
	db.Save(DBFile)
	pn.Save(PNFile)
	writeGob(LUFile, db.LastUpdate)
	fmt.Println("Database updated.")
}
