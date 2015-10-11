package helper

import (
	"archive/zip"
	"encoding/json"
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
	if ze != nil {
		panic(ze)
	}
	defer z.Close()
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
	if oze != nil {
		panic(oze)
	}
	defer ozf.Close()
	jd := json.NewDecoder(ozf)
	var jsf map[string]interface{}
	if je := jd.Decode(jsf); je != nil {
		panic(je)
	}
	return jsf
}

func NewModHelper(filename string) ModHelper {
	return &mcmodinfohelper{filename}
}
