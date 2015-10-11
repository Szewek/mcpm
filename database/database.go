package database

type (
	Database interface {
		Read()
		Update()
		Packages() PkgList
		Files() FileList
	}
	database struct {
		p PkgList
		f FileList
	}
)

func (db *database) Read() {}
func (db *database) Write() {}
func (db *database) Packages() PkgList {
	return nil
}
func (db *database) Files() FileList {
	return nil
}

func NewDatabase() Database {
	return &database{}
}