package util_test

import (
	"testing"

	"github.com/Szewek/mcpm/util"
)

type (
	TestCaseFileSize struct {
		input  int64
		wanted string
	}
)

func TestFileSize(t *testing.T) {
	cases := []TestCaseFileSize{
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
