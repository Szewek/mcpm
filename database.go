package main

import "encoding/gob"

type (
	_FileElement struct {
		ID      int
		Version string
	}
	_DataElement struct {
		ID                         int
		Type                       _PackageType
		PkgName, Name, Description string
		Authors                    []string
		Versions                   []_FileElement
	}
	_Database map[int]_DataElement
)

func (db *_Database) Save(file string) error {
	return writeGobGzip(file, db)
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
