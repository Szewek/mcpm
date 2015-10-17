package helper

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Szewek/mcpm/database"
	"github.com/Szewek/mcpm/util"
)

type (
	ModPackHelper interface {
		Unpack()
	}
	packages struct {
		PID int `json:"projectID"`
		FID int `json:"fileID"`
	}
	modloaderinfo struct {
		ID      string `json:"id"`
		Primary bool   `json:"primary"`
	}
	minecraftinfo struct {
		Version    string          `json:"version"`
		Libraries  string          `json:"libraries"`
		ModLoaders []modloaderinfo `json:"modLoaders"`
	}
	modpackmanifest struct {
		Version        int           `json:"manifestVersion"`
		ModPackVersion string        `json:"version"`
		Files          []packages    `json:"files"`
		Overrides      string        `json:"overrides"`
		Minecraft      minecraftinfo `json:"minecraft"`
	}
	modpackhelper struct {
		filename   string
		filledinfo bool
		info       *modpackmanifest
	}
)

func (mph *modpackhelper) readinfo() {
	if mph.info != nil {
		return
	}
	mph.info = &modpackmanifest{}

	z, ze := zip.OpenReader(mph.filename)
	util.Must(ze)
	defer util.MustClose(z)

	var zf *zip.File
	for i := 0; i < len(z.File); i++ {
		if z.File[i].Name == "manifest.json" {
			zf = z.File[i]
			break
		}
	}
	if zf == nil {
		return
	}

	ozf, oze := zf.Open()
	util.Must(oze)
	defer util.MustClose(ozf)

	jd := json.NewDecoder(ozf)
	util.Must(jd.Decode(mph.info))
	mph.filledinfo = true
}
func (mph *modpackhelper) Unpack() {
	mph.readinfo()
	if mph.filledinfo {
		db := database.GetDatabase()
		for i := 0; i < len(mph.info.Files); i++ {
			fid := mph.info.Files[i].FID
			pkg := db.Packages().Get(mph.info.Files[i].PID)

			if pkg == nil {
				fmt.Printf("Package with ID %d is missing", mph.info.Files[i].PID)
				continue
			}

			fn, pr, hte := util.DownloadPackage(pkg.Type, pkg.ID, pkg.Name, fid)
			util.Must(hte)
			defer util.MustClose(pr)

			pkgo := util.GetPackageOptions(pkg.Type)
			util.Must(util.MkDirIfNotExist(pkgo.Dir))

			sav := fmt.Sprintf("%s/%s", pkgo.Dir, fn)
			f, fe := os.Create(sav)
			util.Must(fe)
			defer util.MustClose(f)

			_, ce := io.Copy(f, pr)
			util.Must(ce)
		}

		z, ze := zip.OpenReader(mph.filename)
		util.Must(ze)
		defer util.MustClose(z)

		for i := 0; i < len(z.File); i++ {
			ov := mph.info.Overrides + "/"
			ix := strings.Index(z.File[i].Name, ov)
			if ix != -1 {
				fln := z.File[i].Name[ix+len(ov):]
				if z.File[i].FileInfo().IsDir() {
					util.MkDirIfNotExist(fln)
				} else {
					zf, zfe := z.File[i].Open()
					util.Must(zfe)

					pr := util.NewProgressReader(zf, z.File[i].UncompressedSize64, fmt.Sprintf("Unpacking %#v...", fln))
					defer util.MustClose(pr)

					f, fe := os.Create(fln)
					util.Must(fe)
					defer util.MustClose(f)

					_, ce := io.Copy(f, pr)
					util.Must(ce)
				}
			}
		}
	}
}

func NewModPackHelper(filename string) ModPackHelper {
	return &modpackhelper{filename, false, nil}
}
