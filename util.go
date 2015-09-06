package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	fsizes        = "kMGTP"
	errLog        = "AN ERROR OCURRED!\n\n=== Go representation of that error ===\n%#v\n\n=== Error message ===\n%s\n-----\n"
	loadBarEmpty  = "          "
	loadBarFilled = "=========="
)

type (
	_ProgressReader struct {
		r        io.Reader
		cnt, tot int64
	}
)

var (
	pkgURLDirs = map[_PackageType]string{
		type_Mod:          "mc-mods",
		type_ModPack:      "modpacks",
		type_ResourcePack: "texture-packs",
		type_WorldSave:    "worlds",
	}
)

func (pr *_ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.cnt += int64(n)
	ld := float64(pr.cnt) / float64(pr.tot)
	ldi := int(ld * 10.0)
	var f string
	if f = "\r"; ld >= 1.0 {
		f = "\n"
	}
	fmt.Printf("[%s%s] %.2f%%    %s", loadBarFilled[0:ldi], loadBarEmpty[ldi:10], ld*100.0, f)
	return
}
func (pr *_ProgressReader) Close() error {
	if rc := pr.r.(io.Closer); rc != nil {
		return rc.Close()
	}
	return nil
}
func newProgressReader(r io.Reader, tot int64) *_ProgressReader {
	return &_ProgressReader{r, 0, tot}
}

func must(e error) {
	if e != nil {
		fmt.Printf(errLog, e, e.Error())
		panic(e)
	}
}
func homePath(file string) string {
	return fmt.Sprint(homeDir, "/", file)
}
func mkDirIfNotExist(dir string) error {
	if _, de := os.Stat(dir); os.IsNotExist(de) {
		return os.MkdirAll(dir, 511)
	}
	return nil
}
func fileSize(size int64) string {
	var s float64
	i := -1
	for s = float64(size); s > 1024 && i < len(fsizes); s = s / 1024 {
		i++
	}
	if i < 0 {
		return fmt.Sprintf("%.3f B", s)
	}
	return fmt.Sprintf("%.3f %cB", s, fsizes[i])
}
func findJava() error {
	cmd := exec.Command("java", "-version")
	var out bytes.Buffer
	cmd.Stderr = &out
	cer := cmd.Run()
	if cer != nil {
		return cer
	}
	s := out.String()
	sst, sen := strings.Index(s, "version"), strings.IndexRune(s, '\n')
	if sst < 0 || sen < sst+8 {
		fmt.Println("Java found but can't get version")
	}
	fmt.Printf("JAVA: %s\n", strings.TrimSpace(s[sst+8:sen]))
	return nil
}
func findGo() error {
	cmd := exec.Command("go", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cer := cmd.Run()
	if cer != nil {
		return cer
	}
	s := out.String()
	sst, sen := strings.Index(s, "version"), strings.IndexRune(s, '\n')
	if sst < 0 || sen < sst+8 {
		fmt.Println("Go found but can't get version")
	}
	fmt.Printf("GO: %s\n", strings.TrimSpace(s[sst+8:sen]))
	return nil
}
func downloadPackage(dt *_DataElement, fid int) (string, *_ProgressReader, error) {
	us := fmt.Sprintf("http://minecraft.curseforge.com/%s/%d-%s/files/", pkgURLDirs[dt.Type], dt.ID, dt.PkgName)
	var us2 string
	if us2 = "latest"; fid != -1 {
		us2 = fmt.Sprintf("%d/download", fid)
	}
	download := fmt.Sprint(us, us2)
	if verbose {
		fmt.Printf("Checking URL %#v\n", download)
	}
	ht, hte := http.Get(download)
	if hte != nil {
		return "", nil, hte
	}
	if verbose {
		fmt.Printf("Found file URL: %#v\n", ht.Request.URL.String())
	}
	fname := ht.Request.URL.Path
	fname = fname[strings.LastIndex(fname, "/")+1:]
	pr := newProgressReader(ht.Body, ht.ContentLength)
	return fname, pr, nil
}
func readGob(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_RDONLY, 0)
	if fe != nil {
		return fe
	}
	return gob.NewDecoder(f).Decode(v)
}
func readGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_RDONLY, 0)
	if fe != nil {
		return fe
	}
	gz, ge := gzip.NewReader(f)
	defer gz.Close()
	if ge != nil {
		return ge
	}
	gb := gob.NewDecoder(gz)
	return gb.Decode(v)
}
func writeGob(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
	if fe != nil {
		return fe
	}
	return gob.NewEncoder(f).Encode(v)
}
func writeGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
	if fe != nil {
		return fe
	}
	gz := gzip.NewWriter(f)
	defer gz.Close()
	gb := gob.NewEncoder(gz)
	return gb.Encode(v)
}
