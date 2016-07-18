package mcpmdb

import (
	"encoding/binary"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Szewek/mcpm/util"
)

const fileName = "v05.mcpmdb"

type (
	mcpmdbImpl struct {
		fn string
		o  bool
		r  map[string]*MCPMPackage
		a  []*MCPMPackage
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
	a:  make([]*MCPMPackage, 0, 4),
}

// GetPackage finds a package by a given name.
// If package doesn't exist in database, then downloads and saves information about new package.
func GetPackage(name string) *MCPMPackage {
	if !db.o {
		open(db.fn)
	}
	if pkg, ok := db.r[name]; ok {
		return pkg
	}
	ht, hte := get(util.DirPathJoin("http://minecraft.curseforge.com/projects", name))
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
	db.r[pn] = pkg
	db.r[idn] = pkg
	db.a = append(db.a, pkg)
	return pkg
}
func open(fname string) {
	db.o = true
	f, fe := os.OpenFile(fname, os.O_RDONLY, 0)
	defer f.Close()
	if os.IsNotExist(fe) {
		return
	}
	util.Must(fe)
	be := binary.BigEndian
	fpb := make([]byte, 20)
	var rp *MCPMPackage
	for {
		_, rer := f.Read(fpb)
		if rer != nil {
			break
		}
		rp = new(MCPMPackage)
		rp.id = be.Uint64(fpb)
		ln := be.Uint32(fpb[8:])
		lt := be.Uint32(fpb[12:])
		lp := be.Uint32(fpb[16:])
		xb := make([]byte, ln+lt+lp)
		_, rer = f.Read(xb)
		util.Must(rer)
		rp.name = string(xb[:ln])
		rp.title = string(xb[ln : ln+lt])
		rp.ptype = string(xb[ln+lt:])
		db.r[rp.name] = rp
		db.r[strconv.FormatUint(rp.id, 10)] = rp
	}
}

// Close checks if a database needs to be updated before application finishes.
func Close() {
	f, fe := os.OpenFile(db.fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer util.MustClose(f)
	util.Must(fe)
	for _, v := range db.a {
		ln := len(v.name)
		lt := len(v.title)
		lp := len(v.ptype)
		tot := 20 + ln + lt + lp
		b := make([]byte, tot)
		be := binary.BigEndian
		bp := 20
		be.PutUint64(b, v.id)
		be.PutUint32(b[8:], uint32(ln))
		be.PutUint32(b[12:], uint32(lt))
		be.PutUint32(b[16:], uint32(lp))
		bp += copy(b[bp:], v.name)
		bp += copy(b[bp:], v.title)
		copy(b[bp:], v.ptype)
		_, wer := f.Write(b)
		util.Must(wer)
	}
}
