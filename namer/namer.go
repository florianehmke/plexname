package namer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type Args struct {
	Path      string
	Overrides parser.Result
}

type Namer struct {
	args Args

	tmdb *tmdb.Service
	tvdb *tvdb.Service
	fs   fs.FileSystem

	files map[string]fileInfo
}

func New(args Args, tmdb *tmdb.Service, tvdb *tvdb.Service, fs fs.FileSystem) *Namer {
	return &Namer{
		args:  args,
		tmdb:  tmdb,
		tvdb:  tvdb,
		fs:    fs,
		files: map[string]fileInfo{},
	}
}

func (pn *Namer) Run() error {
	if err := pn.collectFiles(); err != nil {
		return err
	}
	return nil
}

type fileInfo struct {
	relativePath string
	info         os.FileInfo
}

func (pn *Namer) collectFiles() error {
	if err := filepath.Walk(pn.args.Path, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			pn.files[path] = fileInfo{
				relativePath: strings.TrimLeft(path, pn.args.Path+string(os.PathSeparator)),
				info:         node,
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}
