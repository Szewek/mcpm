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

func MkDirIfNotExist(dir string) error {
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		return os.MkdirAll(dir, 511)
	}
	return nil
}
