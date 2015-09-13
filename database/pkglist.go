package database

type (
	PkgElement struct {
		ID, Type          int
		Name, Description string
		Authors           []string
	}
	PkgList interface {
		Get(pid int) *PkgElement
	}
)
