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
	pkg_WorldSave    _PackageType = 1 // World save
	pkg_ResourcePack _PackageType = 3 // Resource pack
	pkg_ModPack      _PackageType = 5 // Mod pack
	pkg_Mod          _PackageType = 6 // Mod
)

func init() {
	gob.Register(_DataElement{})
	gob.Register(_Database{})
	gob.Register(_PackageNames{})
}
