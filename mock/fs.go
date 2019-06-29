package mock

import "github.com/florianehmke/plexname/fs"

type RenameFn func(oldPath string, newPath string) error
type MkdirAllFn func(path string) error

func NewMockFS(renameFn RenameFn, mkdirAllFn MkdirAllFn) fs.FileSystem {
	return mockFS{
		renameFn:   renameFn,
		mkdirAllFn: mkdirAllFn,
	}
}

type mockFS struct {
	renameFn   func(string, string) error
	mkdirAllFn func(string) error
}

func (fs mockFS) Rename(oldpath, newpath string) error {
	return fs.renameFn(oldpath, newpath)
}

func (fs mockFS) MkdirAll(path string) error {
	return fs.mkdirAllFn(path)
}
