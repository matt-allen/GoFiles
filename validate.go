package fs

import (
	"os"
	"regexp"
)

func isValidFileName(s string) bool {
	m, err := regexp.MatchString("^[^<>:;,?\"*|/]+$", s)
	return err == nil && m
}

func doesFileExist(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func isValidFolderPath(s string) bool {
	return true
}
