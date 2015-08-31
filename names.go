package main

import (
	"encoding/gob"
)

type PackageNames map[string]int

func (pn *PackageNames) FindNames(id int) []string {
	l := []string{}
	for s, i := range *pn {
		if id == i {
			l = append(l, s)
		}
	}
	return l
}
func (pn *PackageNames) Save(file string) {
	writeGobGzip(file, pn)
}
func readPackageNamesFromFile(file string) (*PackageNames, error) {
	var pn PackageNames
	er := readGobGzip(file, &pn)
	return &pn, er
}
func init() {
	gob.Register(PackageNames{})
}
