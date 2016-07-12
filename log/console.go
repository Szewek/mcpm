package log

import (
	"fmt"
	"io"
	"time"
)

type consolelogger byte

var ConsoleLogger Logger = consolelogger(0)

func (cl consolelogger) log(w io.Writer, lt logtype, a ...interface{}) {
	fmt.Fprintf(w, "%s | %s: ", time.Now().Format(time.RFC3339), lt)
	fmt.Fprintln(w, a...)
}
func (cl consolelogger) logf(w io.Writer, lt logtype, format string, a ...interface{}) {
	fmt.Fprintf(w, "%s | %s: ", time.Now().Format(time.RFC3339), lt)
	fmt.Fprintf(w, format, a...)
	fmt.Fprintln(w)
}
