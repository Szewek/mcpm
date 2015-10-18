package util

import "fmt"

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
