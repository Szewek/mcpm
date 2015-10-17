package helper

import (
	"archive/zip"
	"encoding/json"

	"github.com/Szewek/mcpm/util"
)

type (
	ModHelper interface {
		ReadContents() map[string]interface{}
	}
	mcmodinfohelper struct {
		filename string
	}
)

func (mh *mcmodinfohelper) ReadContents() map[string]interface{} {
	z, ze := zip.OpenReader(mh.filename)
	util.Must(ze)
	defer util.MustClose(z)
	var zf *zip.File
	for i := 0; i < len(z.File); i++ {
		if z.File[i].Name == "mcmod.info" {
			zf = z.File[i]
			break
		}
	}
	if zf == nil {
		return nil
	}
	ozf, oze := zf.Open()
	util.Must(oze)
	defer util.MustClose(ozf)
	jd := json.NewDecoder(ozf)
	var jsf map[string]interface{}
	util.Must(jd.Decode(jsf))
	return jsf
}

func NewModHelper(filename string) ModHelper {
	return &mcmodinfohelper{filename}
}
