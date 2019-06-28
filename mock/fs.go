package mock

import "github.com/florianehmke/plexname/fs"

func NewMockFS(err error) fs.FileSystem {
	return mockFS{err: err}
}

type mockFS struct {
	err error
}

func (fs mockFS) Rename(oldpath, newpath string) error {
	return fs.err
}

func (fs mockFS) MkdirAll(path string) error {
	return fs.err
}
