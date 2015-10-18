package util

import (
	"fmt"
	"io"
)

const (
	errLog = "\nAN ERROR OCURRED!\n\n=== Go representation of that error ===\n%#v\n\n=== Error message ===\n%s\n-----\n"
)

// Must handles various errors and prints their Go and string values.
func Must(e error) {
	if e != nil {
		fmt.Printf(errLog, e, e.Error())
		panic(e)
	}
}

// MustClose handles io.Closer errors when deferred.
func MustClose(cl io.Closer) {
	Must(cl.Close())
}
