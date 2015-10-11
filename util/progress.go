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
	ProgressReader struct {
		r    io.Reader
		c, t int64
	}
)

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.c += int64(n)
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
	fmt.Printf("[%s%s] %.2f%%    %s", loadBarFilled[:ldi], loadBarEmpty[ldi:], ld*100.0, f)
	return
}
func (pr *ProgressReader) Close() error {
	if rc := pr.r.(io.Closer); rc != nil {
		return rc.Close()
	}
	return nil
}

func NewProgressReader(r io.Reader, total int64) *ProgressReader {
	return &ProgressReader{r, 0, total}
}
