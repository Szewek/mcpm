package database

type (
	// FileElement contains information about file available on Curse server.
	FileElement struct {
		ID, PkgID       int    // File ID, Package ID
		Name, MCVersion string // File name, Minecraft version
	}
	// FileList interface supplies a method of getting file information
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
