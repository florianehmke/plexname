package fs

import (
	"io"
	"os"
)

type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
}

func NewFileSystem() FileSystem {
	return osFS{}
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (File, error) {
	return os.Open(name)
}

func (osFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
