package util

import (
	"fmt"
	"io"
	"time"
)

const (
	loadBarLen = 40
	updateDur  = time.Duration(200 * time.Millisecond)
)

type (
	// ProgressReader is an io.ReadCloser which outputs progress in a terminal.
	ProgressReader struct {
		r       io.Reader
		l, c, t uint64
		st, n   time.Time
		intro   string
		b       []byte
	}
)

// Read reads data and updates progress.
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	if pr.c == 0 {
		pr.st = time.Now()
		fmt.Println(pr.intro)
	} else if pr.c > pr.t {
		return
	}
	n, err = pr.r.Read(p)
	if n == 0 {
		return
	}
	tn := time.Now()
	pr.c += uint64(n)
	ld := float64(pr.c) / float64(pr.t)
	var f byte
	if f = '\r'; ld >= 1.0 {
		f = '\n'
	}
	dur := tn.Sub(pr.n)
	if ld < 1.0 && dur < updateDur {
		return
	}
	pr.n = tn
	ldi := int(ld * loadBarLen)
	if ldi > loadBarLen {
		ldi = loadBarLen
	} else if ldi < 0 {
		ldi = 0
	}
	for i := 0; i < ldi; i++ {
		pr.b[i] = '='
	}
	dur = tn.Sub(pr.st)
	ds, db := FileSizeNum(float64(pr.c) / (float64(dur) / float64(time.Second)))
	fmt.Printf(" [%s] %6.2f%% %7.2f %cB/s%c", pr.b, ld*100.0, ds, db, f)
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
	pr := &ProgressReader{r: r, t: total, intro: intro, b: make([]byte, loadBarLen)}
	for i := range pr.b {
		pr.b[i] = ' '
	}
	return pr
}
