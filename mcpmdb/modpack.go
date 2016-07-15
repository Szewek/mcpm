package mcpmdb

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Szewek/mcpm/util"
)

type (
	// ModPackHelper provides downloading mods and unpacking files.
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
	buf := make([]byte, 32*1024)
	if mph.filledinfo {
		for i := 0; i < len(mph.info.Files); i++ {
			pid := strconv.FormatInt(int64(mph.info.Files[i].PID), 10)
			fid := mph.info.Files[i].FID
			pkg := GetPackage(pid)
			if pkg == nil {
				fmt.Printf("Package with ID %s is missing!\n", pid)
				continue
			}
			pkg.DownloadFileWithID(strconv.FormatInt(int64(fid), 10), buf)
		}

		z, ze := zip.OpenReader(mph.filename)
		util.Must(ze)
		defer util.MustClose(z)

		ov := mph.info.Overrides + "/"
		for i := 0; i < len(z.File); i++ {
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

					_, ce := io.CopyBuffer(f, pr, buf)
					util.Must(ce)
				}
			}
		}
	}
}

// NewModPackHelper creates a helper for mod pack with a specified file name.
func NewModPackHelper(filename string) ModPackHelper {
	return &modpackhelper{filename, false, nil}
}
