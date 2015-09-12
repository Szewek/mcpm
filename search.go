package main

import (
	"fmt"
	"regexp"
)

type (
	IDArray []int
)

var pkgTypeName = map[_PackageType]string{
	pkg_Mod:          "Mod",
	pkg_ModPack:      "Modpack",
	pkg_ResourcePack: "Resource pack",
	pkg_WorldSave:    "World save",
}

func (ia *IDArray) Size() int {
	return len(*ia)
}
func (ia *IDArray) Get(idx int) int {
	return (*ia)[idx]
}
func (ia *IDArray) Contains(i int) bool {
	for j := 0; j < len(*ia); j++ {
		if i == (*ia)[j] {
			return true
		}
	}
	return false
}
func (ia *IDArray) Add(i int) {
	if !ia.Contains(i) {
		(*ia) = append(*ia, i)
	}
}

func searchPackage() {
	val := flagset.Arg(0)
	rgx := regexp.MustCompile(val)
	fmt.Printf("Searching %#v...\n", val)
	ids := new(IDArray)
	for p, i := range *_PKGS {
		idx := rgx.FindStringIndex(p)
		if idx != nil {
			ids.Add(i)
		}
	}
	for i, dt := range *_DBASE {
		idx := rgx.FindStringIndex(dt.Name)
		if idx != nil {
			ids.Add(i)
		}
		idx = rgx.FindStringIndex(dt.Description)
		if idx != nil {
			ids.Add(i)
		}
	}
	for i := 0; i < ids.Size(); i++ {
		data := (*_DBASE)[ids.Get(i)]
		fmt.Printf("\n%s [%s] – %s\n %s\n", data.Name, data.PkgName, pkgTypeName[data.Type], data.Description)
	}
}
