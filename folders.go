package fs

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func getFolderFromPath(p string) string {
	split := strings.Split(p, "/")
	return strings.Join(split[:len(split)-1], "/")
}

func getFileNameFromPath(p string) string {
	split := strings.Split(p, "/")
	return strings.Join(split[len(split)-1:], "/")
}

// List folders within a given directory
func listFolder(f string) []string {
	all := []string{}
	_ = filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				all = append(all, substring(path, len(f)+1))
			}
		}
		return err
	})
	return all
}

func substring(s string, i int) string {
	r := []rune(s)
	return string(r[i:len(r)])
}

// Move a file from one location to another
func moveFile(s, d string) error {
	// Need to make sure the full hierarchy of the destination exists
	_ = os.MkdirAll(d, 0755)
	// Firstly, try renaming it. This will fail if the files are on different partitions.
	err := os.Rename(s, d)
	// If that fails then we need to make a copy of the file, then delete the existing one.
	if err != nil {
		// When copying, it could take a very long time!
		in, err := os.Open(s)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(d)
		if err != nil {
			return err
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()
		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		err = out.Sync()
		// And then delete the existing file
		return deleteFile(s)
	}
	return err
}

func deleteFile(p string) error {
	return os.Remove(p)
}

func fileSize(p string) (uint64, error) {
	in, err := os.Open(p)
	if err != nil {
		return 0, err
	}
	fi, err := in.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(fi.Size()), err
}
