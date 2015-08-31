package main

import (
	"encoding/gob"
)

type _PackageNames map[string]int

func (pn *_PackageNames) FindNames(id int) []string {
	l := []string{}
	for s, i := range *pn {
		if id == i {
			l = append(l, s)
		}
	}
	return l
}
func (pn *_PackageNames) Save(file string) {
	writeGobGzip(file, pn)
}
func readPackageNamesFromFile(file string) (*_PackageNames, error) {
	var pn _PackageNames
	er := readGobGzip(file, &pn)
	return &pn, er
}
func init() {
	gob.Register(_PackageNames{})
}
