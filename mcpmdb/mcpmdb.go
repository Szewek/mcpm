package mcpmdb

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Szewek/mcpm/util"
)

const fileName = "v05.mcpmdb"

type (
	mcpmdbImpl struct {
		fn     string
		opened bool
		r, w   map[string]*MCPMPackage
	}
	// MCPMPackage represents a Minecraft package to download
	MCPMPackage struct {
		id                 uint64
		name, title, ptype string
	}
)

var db = &mcpmdbImpl{
	fn: util.GetHomeDir() + fileName,
	r:  map[string]*MCPMPackage{},
	w:  map[string]*MCPMPackage{},
}

// GetPackage finds a package by a given name.
// If package doesn't exist in database, then downloads and saves information about new package.
func GetPackage(name string) *MCPMPackage {
	if !db.opened {
		open(db.fn)
	}
	if pkg, ok := db.r[name]; ok {
		return pkg
	}
	if pkg, ok := db.w[name]; ok {
		return pkg
	}
	ur := fmt.Sprintf(util.CurseForgeURL, name)
	ht, hte := http.Get(ur)
	if hte != nil {
		return nil
	}
	url := ht.Request.URL
	pn := url.Path[strings.LastIndexByte(url.Path, '/')+1:]
	doc, doce := goquery.NewDocumentFromResponse(ht)
	if doce != nil {
		return nil
	}
	tit := doc.Find("h1.project-title a .overflow-tip").Text()
	oa := doc.Find(".RootGameCategory a").Eq(0)
	ida := doc.Find(".view-on-curse a").Eq(0)
	cat := oa.AttrOr("href", " ")
	id := ida.AttrOr("href", "")
	idn := id[strings.LastIndexByte(id, '/')+1:]
	nid, ner := strconv.ParseUint(idn, 10, 64)
	if ner != nil {
		return nil
	}
	pkg := &MCPMPackage{nid, pn, tit, cat[1:]}
	db.w[pn] = pkg
	db.w[idn] = pkg
	return pkg
}
func open(fname string) {
	db.opened = true
	f, fe := os.OpenFile(fname, os.O_RDONLY, 0)
	defer f.Close()
	if fe != nil {
		return
	}
	gb := gob.NewDecoder(f)
	var rp *MCPMPackage
	var eof error
	for {
		rp = new(MCPMPackage)
		eof = gb.Decode(rp)
		if eof != nil {
			break
		}
		db.r[rp.name] = rp
		db.r[strconv.FormatUint(rp.id, 10)] = rp
	}
}

// Close checks if a database needs to be updated before application finishes.
func Close() error {
	f, fe := os.OpenFile(db.fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if fe != nil {
		return fe
	}
	a := make([]*MCPMPackage, len(db.w))
	ap := 0
LOOP1:
	for _, v := range db.w {
		for i := 0; i < ap; i++ {
			if a[i] == v {
				continue LOOP1
			}
		}
		a[ap] = v
		ap++
	}
	a = a[:ap]
	gb := gob.NewEncoder(f)
	for _, v := range a {
		gbe := gb.Encode(v)
		if gbe != nil {
			return gbe
		}
	}
	return nil
}

func init() {
	gob.Register(&MCPMPackage{})
}
