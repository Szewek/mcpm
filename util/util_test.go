package util_test

import (
	"testing"

	"github.com/Szewek/mcpm/util"
)

type (
	testCaseFileSize struct {
		input  int64
		wanted string
	}
	noReader byte
)

const bufSize = 4096

func (nr *noReader) Read(b []byte) (int, error) {
	return len(b), nil
}

func TestFileSize(t *testing.T) {
	cases := []testCaseFileSize{
		{2048, "2.00 kB"},
		{1024, "1.00 kB"},
		{64, "64 B"},
	}
	for i := 0; i < len(cases); i++ {
		result := util.FileSize(cases[i].input)
		if result != cases[i].wanted {
			t.Errorf("Expected %#v, got %#v!", cases[i].wanted, result)
		} else {
			t.Logf("Got %#v from %d", result, cases[i].input)
		}
	}
}
func BenchmarkFileSize(b *testing.B) {
	for t := 0; t < b.N; t++ {
		util.FileSize(65536)
	}
}

func BenchmarkProgress(b *testing.B) {
	tot := uint64(b.N * bufSize)
	buf := make([]byte, bufSize)
	prog := util.NewProgressReader(new(noReader), tot, "")
	b.SetBytes(bufSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prog.Read(buf)
	}
}
