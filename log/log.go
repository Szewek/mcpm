package log

import (
	"io"
	"os"
)

type (
	logtype int
	Logger  interface {
		log(io.Writer, logtype, ...interface{})
		logf(io.Writer, logtype, string, ...interface{})
	}
	LogWriter struct {
		w io.Writer
		l Logger
	}
)

const (
	info logtype = iota
	warning
	err
)

var (
	lognames []string = []string{
		"INFO",
		"WARNING",
		"ERROR",
	}
)

func (lt logtype) String() string {
	return lognames[lt]
}

func (lw *LogWriter) Info(a ...interface{}) {
	lw.l.log(lw.w, info, a)
}
func (lw *LogWriter) InfoF(format string, a ...interface{}) {
	lw.l.logf(lw.w, info, format, a)
}
func (lw *LogWriter) Warn(a ...interface{}) {
	lw.l.log(lw.w, warning, a)
}
func (lw *LogWriter) WarnF(format string, a ...interface{}) {
	lw.l.logf(lw.w, warning, format, a)
}
func (lw *LogWriter) Err(a ...interface{}) {
	lw.l.log(lw.w, err, a)
}
func (lw *LogWriter) ErrF(format string, a ...interface{}) {
	lw.l.logf(lw.w, err, format, a)
}

func NewLogWriter(l Logger) *LogWriter {
	return &LogWriter{os.Stdout, l}
}

func init() {
	NewLogWriter(ConsoleLogger).Info("Started logging")
}
