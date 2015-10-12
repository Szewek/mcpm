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
		Each(fn func(int, *PkgElement))
	}
	pkglist struct {
		Names map[string]int
		Pkgs  map[int]PkgElement
	}
)

func (pl *pkglist) Get(pid int) *PkgElement {
	if el, ok := pl.Pkgs[pid]; ok {
		return &el
	}
	return nil
}
func (pl *pkglist) GetByName(pname string) *PkgElement {
	if id, ok := pl.Names[pname]; ok {
		return pl.Get(id)
	}
	return nil
}
func (pl *pkglist) Each(fn func(int, *PkgElement)) {
	for k, v := range pl.Pkgs {
		fn(k, &v)
	}
}
