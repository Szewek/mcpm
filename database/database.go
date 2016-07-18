// Package database helps maintaining Minecraft resources available at Curse CDN
package database

import (
	"encoding/gob"

	"github.com/Szewek/mcpm/util"
)

const dbFile = "v040.mcpmdb"

type (
	// Database interface allows to store lists of files and packages locally or remotely.
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
	db    = &database{}
	dbset = false
)

// GetDatabase function returns a database stored in MCPM local directory.
// If database doesn't exist it gets updated.
func GetDatabase() Database {
	if !dbset {
		util.ReadGobGzip(dbFile, db)
		dbset = true
	}
	return db
}

func init() {
	gob.Register(&database{})
	gob.Register(&pkglist{})
	gob.Register(&filelist{})
}
