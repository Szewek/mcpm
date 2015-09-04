package main

import (
	"fmt"
	"regexp"
)

var pkgTypeName = map[_PackageType]string{
	type_Mod:          "Mod",
	type_ModPack:      "Modpack",
	type_ResourcePack: "Resource pack",
	type_WorldSave:    "World save",
}

func searchPackage() {
	val := flagset.Arg(0)
	pn, pne := readPackageNamesFromFile(homePath(pnFile))
	must(pne)
	db, dbe := readDatabaseFromFile(homePath(dbFile))
	must(dbe)
	rgx := regexp.MustCompile(val)
	fmt.Printf("Searching %#v...\n", val)
	ids := []int{}
	for p, i := range *pn {
		idx := rgx.FindStringIndex(p)
		if idx != nil {
			ids = append(ids, i)
		}
	}
	for i, dt := range *db {
		idx := rgx.FindStringIndex(dt.Name)
		if idx != nil {
			ids = append(ids, i)
		}
		idx = rgx.FindStringIndex(dt.Description)
		if idx != nil {
			ids = append(ids, i)
		}
	}
	for i := 0; i < len(ids); i++ {
		data := (*db)[ids[i]]
		fmt.Printf("\n%s [%s] – %s\n %s\n", data.Name, data.PkgName, pkgTypeName[data.Type], data.Description)
	}
}
