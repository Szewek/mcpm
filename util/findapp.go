package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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
func FindGo() error {
	cmd := exec.Command("go", "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cer := cmd.Run()
	if cer != nil {
		return cer
	}
	s := out.String()
	sst, sen := strings.Index(s, "version"), strings.IndexRune(s, '\n')
	if sst < 0 || sen < sst+8 {
		fmt.Println("Go found but can't get version")
	}
	fmt.Printf("GO: %s\n", strings.TrimSpace(s[sst+8:sen]))
	return nil
}
