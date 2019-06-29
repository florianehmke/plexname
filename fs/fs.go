package fs

import (
	"os"
)

type FileSystem interface {
	Rename(oldpath, newpath string) error
	MkdirAll(path string) error
}

func NewFileSystem() FileSystem {
	return osFS{}
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (osFS) MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
