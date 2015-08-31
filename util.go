package main

import (
	"compress/gzip"
	"encoding/gob"
	"os"
)

func must(e error) {
	if e != nil {
		panic(e)
	}
}
func readGob(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_RDONLY, 0)
	if fe != nil {
		return fe
	}
	return gob.NewDecoder(f).Decode(v)
}
func readGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_RDONLY, 0)
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
func writeGob(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
	if fe != nil {
		return fe
	}
	return gob.NewEncoder(f).Encode(v)
}
func writeGobGzip(file string, v interface{}) error {
	f, fe := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 438)
	if fe != nil {
		return fe
	}
	gz := gzip.NewWriter(f)
	defer gz.Close()
	gb := gob.NewEncoder(gz)
	return gb.Encode(v)
}
