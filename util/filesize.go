package util

import "fmt"

const fsizes = "kMGTP"

func FileSize(size int64) string {
	var s float64
	i := -1
	for s = float64(size); s > 1024 && i < len(fsizes); s = s / 1024 {
		i++
	}
	if i < 0 {
		return fmt.Sprintf("%.3f B", s)
	}
	return fmt.Sprintf("%.3f %cB", s, fsizes[i])
}
