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
	Must(mkDirIfNotExist())
}
func mkDirIfNotExist() error {
	if _, de := os.Stat(homedir); os.IsNotExist(de) {
		return os.MkdirAll(homedir, 511)
	}
	return nil
}
