//Package util contains varoius utilities used by MCPM
package util

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	fsizes     = "kMGTP"
	byteStep   = 1000
	byteBiStep = 1024
	errLog     = "\nAN ERROR OCURRED!\n\n=== Go representation of that error ===\n%#v\n\n=== Error message ===\n%s\n-----\n"
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

// FileSize returns a formatted data size in bytes, kilobytes etc.
func FileSize(size int64) string {
	var s float64
	i := -1
	for s = float64(size); s >= byteBiStep && i < len(fsizes); s = s / byteBiStep {
		i++
	}
	if i < 0 {
		return fmt.Sprintf("%d B", int(s))
	}
	return fmt.Sprintf("%.2f %cB", s, fsizes[i])
}

// FileSizeNum returns a calculated data size in bytes, kilobytes etc.
//
// Calculations aren't the same as FileSize does. It is commonly used by ProgressReader.
func FileSizeNum(size float64) (float64, byte) {
	var s float64
	i := -1
	for s = size; s >= byteStep && i < len(fsizes); s = s / byteStep {
		i++
	}
	if i < 0 {
		return s, 0
	}
	if i >= len(fsizes) {
		i = len(fsizes) - 1
	}
	return s, fsizes[i]
}

// FindJava checks if Java is installed.
func FindJava() error {
	cmd := exec.Command("java", "-version")
	var out bytes.Buffer
	cmd.Stderr = &out
	cer := cmd.Run()
	if cer != nil {
		return cer
	}
	s := out.String()
	sst, sen := strings.Index(s, "version"), strings.IndexRune(s, '\n')
	if sst < 0 || sen < sst+8 {
		fmt.Println("Java found but can't get version")
	}
	fmt.Printf("JAVA: %s\n", strings.TrimSpace(s[sst+8:sen]))
	return nil
}
