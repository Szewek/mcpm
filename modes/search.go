package modes

import (
	"fmt"
	"regexp"

	"github.com/Szewek/mcpm/database"
)

type (
	idArray []int
)

func (ia *idArray) Size() int {
	return len(*ia)
}
func (ia *idArray) Get(idx int) int {
	return (*ia)[idx]
}
func (ia *idArray) Contains(i int) bool {
	for j := 0; j < len(*ia); j++ {
		if i == (*ia)[j] {
			return true
		}
	}
	return false
}
func (ia *idArray) Add(i int) {
	if !ia.Contains(i) {
		(*ia) = append(*ia, i)
	}
}

var pkgTypeName = map[int]string{
	6: "Mod",
	5: "Modpack",
	3: "Resource pack",
	1: "World save",
}

func search(mo *ModeOptions) {
	v := mo.Args[0]
	db := database.GetDatabase()
	rgx := regexp.MustCompile(v)
	fmt.Printf("Searching %#v...\n", v)
	ids := new(idArray)
	db.Packages().Each(func(i int, pe *database.PkgElement) {
		idx := rgx.FindStringIndex(pe.Name)
		if idx != nil {
			ids.Add(i)
		}
		idx = rgx.FindStringIndex(pe.FullName)
		if idx != nil {
			ids.Add(i)
		}
		idx = rgx.FindStringIndex(pe.Description)
		if idx != nil {
			ids.Add(i)
		}
	})
	for i := 0; i < ids.Size(); i++ {
		data := db.Packages().Get(ids.Get(i))
		fmt.Printf("\n%s [%s] – %s\n %s\n", data.FullName, data.Name, pkgTypeName[data.Type], data.Description)
	}
}

func init() {
	registerMode("search", search)
}
