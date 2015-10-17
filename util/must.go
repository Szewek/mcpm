package util

import (
	"fmt"
	"io"
)

const (
	errLog = "\nAN ERROR OCURRED!\n\n=== Go representation of that error ===\n%#v\n\n=== Error message ===\n%s\n-----\n"
)

func Must(e error) {
	if e != nil {
		fmt.Printf(errLog, e, e.Error())
		panic(e)
	}
}

func MustClose(cl io.Closer) {
	Must(cl.Close())
}
