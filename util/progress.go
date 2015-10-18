package util

import (
	"fmt"
	"io"
)

const (
	loadBarEmpty  = "          "
	loadBarFilled = "=========="
)

type (
	// ProgressReader is an io.ReadCloser which outputs progress in a terminal.
	ProgressReader struct {
		r     io.Reader
		c, t  uint64
		intro string
	}
)

// Read reads data and updates progress.
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	if pr.c == 0 {
		fmt.Println(pr.intro)
	} else if pr.c > pr.t {
		return
	}
	n, err = pr.r.Read(p)
	if n == 0 {
		return
	}
	pr.c += uint64(n)
	ld := float64(pr.c) / float64(pr.t)
	ldi := int(ld * 10.0)
	var f string
	if f = "\r"; ld >= 1.0 {
		f = "\n"
	}
	if ldi > 10 {
		ldi = 10
	} else if ldi < 0 {
		ldi = 0
	}
	fmt.Printf("  [%s%s] %.2f%%    %s", loadBarFilled[:ldi], loadBarEmpty[ldi:], ld*100.0, f)
	return
}

// Close checks if given io.Reader is also an io.Closer.
// If true, it closes.
func (pr *ProgressReader) Close() error {
	if rc := pr.r.(io.Closer); rc != nil {
		return rc.Close()
	}
	return nil
}

// NewProgressReader returns new ProgressReader.
func NewProgressReader(r io.Reader, total uint64, intro string) *ProgressReader {
	return &ProgressReader{r, 0, total, intro}
}
