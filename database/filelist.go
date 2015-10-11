package database

type (
	FileElement struct {
		ID, PkgID       int
		Name, MCVersion string
	}
	FileList interface {
		Get(fid int) *FileElement
	}
	filelist map[int]FileElement
)

func (fl *filelist) Get(fid int) *FileElement {
	if el, ok := (*fl)[fid]; ok {
		return &el
	}
	return nil
}
