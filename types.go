package main

import "encoding/gob"

type (
	_PackageType int
	_DataElement struct {
		ID                         int
		Type                       _PackageType
		PkgName, Name, Description string
		Authors                    []string
	}
	_Database     map[int]_DataElement
	_PackageNames map[string]int
)

const (
	type_WorldSave    _PackageType = 1 // World save
	type_ResourcePack _PackageType = 3 // Resource pack
	type_ModPack      _PackageType = 5 // Mod pack
	type_Mod          _PackageType = 6 // Mod
)

func (db *_Database) Save(file string) error {
	return writeGobGzip(file, db)
}
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

func init() {
	gob.Register(_DataElement{})
	gob.Register(_Database{})
	gob.Register(_PackageNames{})
}
