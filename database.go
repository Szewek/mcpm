package main

import (
	"encoding/gob"
	"time"
)

type (
	FileElement struct {
		Id      int
		Version string
	}
	DataElement struct {
		Id                         int
		Type                       PackageType
		PkgName, Name, Description string
		Authors                    []string
		Versions                   []FileElement
	}
	Database struct {
		Version    int
		LastUpdate time.Time
		Data       map[int]DataElement
	}
)

func (db *Database) Save(file string) error {
	return writeGobGzip(file, db)
}

func createDatabase() *Database {
	return &Database{1, time.Now(), map[int]DataElement{}}
}
func readDatabaseFromFile(file string) (*Database, error) {
	var db Database
	er := readGobGzip(file, &db)
	return &db, er
}

func init() {
	gob.Register(FileElement{})
	gob.Register(DataElement{})
	gob.Register(Database{})
}
