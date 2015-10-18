package util

import (
	"compress/gzip"
	"encoding/gob"
	"os"
)

// ReadGobGzip reads gzipped and encoded with gob data from a file and decodes it.
func ReadGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(GetHomeDir()+file, os.O_RDONLY, 0)
	if fe != nil {
		return fe
	}
	gz, ge := gzip.NewReader(f)
	defer gz.Close()
	if ge != nil {
		return ge
	}
	gb := gob.NewDecoder(gz)
	return gb.Decode(v)
}

// WriteGobGzip writes gzipped and encoded with gob data to a file.
func WriteGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(GetHomeDir()+file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
	if fe != nil {
		return fe
	}
	gz := gzip.NewWriter(f)
	defer gz.Close()
	gb := gob.NewEncoder(gz)
	return gb.Encode(v)
}
