package database

type (
	FileElement struct {
		ID, PkgID       int
		Name, MCVersion string
	}
	FileList interface {
		Get(fid int) *FileElement
	}
)
