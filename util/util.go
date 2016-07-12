package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const fsizes = "kMGTP"

// FileSize returns a formatted data size in bytes, kilobytes etc.
func FileSize(size int64) string {
	var s float64
	i := -1
	for s = float64(size); s >= 1024 && i < len(fsizes); s = s / 1024 {
		i++
	}
	if i < 0 {
		return fmt.Sprintf("%d B", int(s))
	}
	return fmt.Sprintf("%.2f %cB", s, fsizes[i])
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
