package util

import (
	"fmt"
	"os"
	"os/user"
)

var (
	homedir = ""
	homeset = false
)

// GetHomeDir returns default MCPM directory.
// If directory is not set yet, function creates a path to it and checks if it exists.
func GetHomeDir() string {
	if !homeset {
		createHomeDir()
		homeset = true
	}
	return homedir
}

func createHomeDir() {
	u, ue := user.Current()
	Must(ue)
	homedir = fmt.Sprintf("%s/.mcpm/", u.HomeDir)
	Must(MkDirIfNotExist(homedir))
}

// MkDirIfNotExist checks if a directory exists.
// If not, it creates one.
func MkDirIfNotExist(dir string) error {
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		return os.MkdirAll(dir, 511)
	}
	return nil
}
