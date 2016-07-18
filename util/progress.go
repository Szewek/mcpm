package util

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	loadBarLen  = 30
	floatSecond = float64(time.Second)
	updateDur   = time.Duration(250) * time.Millisecond
)

var (
	barStart = []byte(" [\x1B[33;1m")
	barEnd   = []byte("\x1B[39;49m] ")
)

type (
	progressRead struct {
		sync.Mutex
		r    io.Reader
		x    time.Duration
		t, n float64
		b    []byte
		c    chan progress
	}
	progress struct {
		t time.Duration
		n float64
	}
)

func (p *progressRead) listen() {
	p.Lock()
	defer p.Unlock()
	tx := time.NewTicker(updateDur)
	pl := 9
	var db byte
	var ld, ds float64
	var ldi, i int
	var f byte = '\r'
	for x := range p.c {
		if p.n >= p.t {
			continue
		}
		p.n += x.n
		p.x += x.t
		if p.n >= p.t {
			f = '\n'
		} else {
			select {
			case <-tx.C:
				break
			default:
				continue
			}
		}
		ld = p.n / p.t
		if ld > 1 {
			ld = 1
		}
		ldi = 9 + int(ld*loadBarLen)
		if ldi > pl {
			for i := pl; i < ldi; i++ {
				p.b[i] = '='
			}
			pl = ldi
		}
		ds = p.n * floatSecond / float64(p.x)
		db = ' '
		if ds >= byteStep {
			i = 0
			for ds >= byteStep && i < len(fsizes) {
				ds /= byteStep
				i++
			}
			db = fsizes[i-1]
		}
		ld *= 100
		os.Stdout.Write(p.b)
		fmt.Fprintf(os.Stdout, "%6.2f%% %6.2f %cB/s%c", ld, ds, db, f)
	}
	tx.Stop()
}

func (p *progressRead) Read(b []byte) (n int, err error) {
	t := time.Now()
	n, err = p.r.Read(b)
	p.c <- progress{time.Now().Sub(t), float64(n)}
	return
}

func (p *progressRead) Close() error {
	close(p.c)
	p.Lock()
	defer p.Unlock()
	return nil
}

// NewReadProgress creates a progress bar for a reader.
//
// Read progress must be closed immediatelly after reading operation.
func NewReadProgress(r io.Reader, total uint64) io.ReadCloser {
	b := make([]byte, 19+loadBarLen)
	copy(b, barStart)
	for i := 9; i < loadBarLen+9; i++ {
		b[i] = ' '
	}
	copy(b[loadBarLen+9:], barEnd)
	p := &progressRead{
		r: r,
		t: float64(total),
		b: b,
		c: make(chan progress),
	}
	go p.listen()
	return p
}
