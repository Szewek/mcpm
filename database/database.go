package database

import (
	"encoding/gob"
	"fmt"

	"github.com/Szewek/mcpm/util"
)

const dbFile = "v040.mcpmdb"

type (
	Database interface {
		Packages() PkgList
		Files() FileList
	}
	database struct {
		P *pkglist
		F *filelist
	}
)

func (db *database) Packages() PkgList {
	return db.P
}
func (db *database) Files() FileList {
	return db.F
}

var (
	db    *database = &database{}
	dbset           = false
)

func GetDatabase() Database {
	if !dbset {
		if er := util.ReadGobGzip(dbFile, db); er != nil {
			fmt.Printf("Error DB: %#v\n\n%s", er, er)
			UpdateDatabase(false)
		}
	}
	return db
}

func init() {
	gob.Register(&database{})
	gob.Register(&pkglist{})
	gob.Register(&filelist{})
}
