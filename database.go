package main

import (
	"encoding/gob"
	"time"
)

type (
	_FileElement struct {
		ID      int
		Version string
	}
	_DataElement struct {
		ID                         int
		Type                       PackageType
		PkgName, Name, Description string
		Authors                    []string
		Versions                   []_FileElement
	}
	_Database struct {
		Version    int
		LastUpdate time.Time
		Data       map[int]_DataElement
	}
)

func (db *_Database) Save(file string) error {
	return writeGobGzip(file, db)
}

func createDatabase() *_Database {
	return &_Database{1, time.Now(), map[int]_DataElement{}}
}
func readDatabaseFromFile(file string) (*_Database, error) {
	var db _Database
	er := readGobGzip(file, &db)
	return &db, er
}

func init() {
	gob.Register(_FileElement{})
	gob.Register(_DataElement{})
	gob.Register(_Database{})
}
