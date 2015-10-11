package database

const dbFile = "v040.mcpmdb"

type (
	Database interface {
		Packages() PkgList
		Files() FileList
	}
	database struct {
		p *pkglist
		f *filelist
	}
)

func (db *database) Packages() PkgList {
	return db.p
}
func (db *database) Files() FileList {
	return db.f
}

var db *database = nil

func GetDatabase() Database {
	return db
}
