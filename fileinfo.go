package fs

import "fmt"

type FileList struct {
	path  string
	files []string
}

type FileInfo struct {
	FileName string `json:"file_name"`
	Folder   string `json:"folder"`
	FullPath string `json:"full_path"`
	Size     uint64 `json:"file_size"`
}

func ScanDirectory(path string) FileList {
	return FileList{path, listFolder(path)}
}

func (f *FileList) GetList() []string {
	return f.files
}

func (f *FileList) FreeSpace() uint64 {
	s, e := remainingSpace(f.path)
	if e != nil {
		return 0
	}
	return s
}

func (f *FileList) Enrich(lock *FileOperation) []FileInfo {
	enrd := []FileInfo{}
	for _, i := range f.files {
		fullPath := f.prefixPath(i)
		if lock.IsLocked(fullPath) {
			continue
		}
		size, err := fileSize(fullPath)
		if err != nil {
			size = 0
		}
		enrd = append(enrd, FileInfo{
			getFileNameFromPath(i),
			getFolderFromPath(i),
			fullPath,
			size,
		})
	}
	return enrd
}

func (f *FileList) prefixPath(p string) string {
	return fmt.Sprintf("%s/%s", f.path, p)
}
