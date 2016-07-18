package helper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/Szewek/mcpm/util"
)

type (
	// SaveHelper helps with unpacking files to "saves" folder
	SaveHelper interface {
		UnpackAll()
	}
	savehelper struct {
		filename string
	}
)

func (sh *savehelper) getzip() *zip.ReadCloser {
	z, ze := zip.OpenReader(sh.filename)
	util.Must(ze)
	return z
}
func (sh *savehelper) UnpackAll() {
	z := sh.getzip()
	defer util.MustClose(z)
	for i := 0; i < len(z.File); i++ {
		zi := z.File[i]
		fln := "saves/" + zi.Name

		if zi.FileInfo().IsDir() {
			util.Must(util.MkDirIfNotExist(fln))
		} else {
			zf, zfe := zi.Open()
			util.Must(zfe)

			fmt.Printf("Unpacking %#v...\n", zi.Name)
			pr := util.NewReadProgress(zf, zi.UncompressedSize64)

			f, fe := os.Create(fln)
			util.Must(fe)
			defer util.MustClose(f)

			_, ce := io.Copy(f, pr)
			util.Must(ce)
			util.MustClose(zf)
			util.MustClose(pr)
		}
	}
}

// NewSaveHelper creates SaveHelper that reads content from downloaded package
func NewSaveHelper(filename string) SaveHelper {
	return &savehelper{filename}
}
