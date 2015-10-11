package database

type (
	PkgElement struct {
		ID, Type                    int
		Name, FullName, Description string
		Authors                     []string
	}
	PkgList interface {
		Get(pid int) *PkgElement
		GetByName(pname string) *PkgElement
	}
	pkglist struct {
		names map[string]int
		pkgs  map[int]PkgElement
	}
)

func (pl *pkglist) Get(pid int) *PkgElement {
	if el, ok := pl.pkgs[pid]; ok {
		return &el
	}
	return nil
}
func (pl *pkglist) GetByName(pname string) *PkgElement {
	if id, ok := pl.names[pname]; ok {
		return pl.Get(id)
	}
	return nil
}
