package mcpmdb

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Szewek/mcpm/util"
)

type (
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

// UnpackModpack downloads mods and unpacks file contents.
func UnpackModpack(fname string) {
	z, ze := zip.OpenReader(fname)
	util.Must(ze)
	defer util.MustClose(z)
	info := new(modpackmanifest)
	var zf *zip.File
	ov := "overrides/"
	buf := make([]byte, 32*1024)
	for _, zf = range z.File {
		if zf.Name == "manifest.json" {
			ozf, oze := zf.Open()
			util.Must(oze)
			util.Must(json.NewDecoder(ozf).Decode(info))
			util.MustClose(ozf)
			break
		}
	}
	if info != nil {
		ov = info.Overrides + "/"
		for i, xf := range info.Files {
			pid := strconv.FormatInt(int64(xf.PID), 10)
			pkg := GetPackage(pid)
			if pkg == nil {
				fmt.Printf("Package with ID %s is missing!\n", pid)
				continue
			}
			fmt.Printf("%d / %d ", i, len(info.Files))
			pkg.DownloadFileWithID(strconv.FormatInt(int64(xf.FID), 10), buf)
		}
	}
	lov := len(ov)
	for _, zf = range z.File {
		if len(zf.Name) < lov {
			continue
		}
		n := zf.Name[:lov]
		if n == ov {
			n = zf.Name[lov:]
			if n == "" {
				continue
			}
			if zf.FileInfo().IsDir() {
				util.Must(util.MkDirIfNotExist(n))
			} else {
				xf, xe := zf.Open()
				util.Must(xe)
				fmt.Printf("Unpacking %#v...\n", n)
				pr := util.NewReadProgress(xf, zf.UncompressedSize64)
				f, fe := os.Create(n)
				util.Must(fe)
				_, ce := io.CopyBuffer(f, pr, buf)
				util.Must(ce)
				util.MustClose(f)
				util.MustClose(xf)
				pr.Close()
			}
		}
	}
}
