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
