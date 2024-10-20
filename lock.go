package fs

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
)

type FileOperation struct {
	locks map[string]uint64
}

func CreateFileLock() *FileOperation {
	return &FileOperation{make(map[string]uint64)}
}

func (f *FileOperation) Move(from, to string) {
	size, err := fileSize(from)
	if err != nil {
		f.locks[from] = 0
	} else {
		f.locks[from] = size
	}
	_ = moveFile(from, to)
	delete(f.locks, from)
	log.Println(fmt.Sprintf("Moving of %s is complete", from))
}

func (f *FileOperation) IsLocked(p string) bool {
	_, exists := f.locks[p]
	return exists
}

func (f *FileOperation) CanMove(from, to string) error {
	log.Println(fmt.Sprintf("Requesting to move %s to %s", from, to))
	if !doesFileExist(from) {
		return errors.New("the source file does not exist")
	}
	if doesFileExist(to) {
		return errors.New("file already exists with this name at the destination")
	}
	if !isValidFolderPath(from) {
		return errors.New("the source is not a valid file path")
	}
	if !isValidFolderPath(to) {
		return errors.New("the destination is not a valid file path")
	}
	if f.IsLocked(from) {
		return errors.New("this file is already being moved and cannot be changed")
	}
	size, sErr := fileSize(from)
	if sErr != nil {
		return sErr
	}
	log.Println("Checking remaining space at " + getFolderFromPath(to))
	left, lErr := remainingSpace(getFolderFromPath(to))
	if lErr != nil {
		return lErr
	}
	log.Println(fmt.Sprintf("File size: %d, remaining: %d", size, left))
	if size > left {
		return errors.New("not enough space at the destination to move this file")
	}
	// TODO Also check write permissions in destination directory
	return nil
}

func (f *FileOperation) CanDelete(from string) bool {
	return isValidFolderPath(from) && !f.IsLocked(from)
}

func (f *FileOperation) Delete(from string) error {
	return deleteFile(from)
}

func remainingSpace(p string) (uint64, error) {
	var stat unix.Statfs_t
	err := unix.Statfs(p, &stat)
	return stat.Bavail * uint64(stat.Bsize), err
}
