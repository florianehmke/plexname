package namer

import (
	"fmt"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
	"os"
	"path/filepath"
)

type Args struct {
	Path      string
	Overrides parser.Result
}

type Namer struct {
	args Args

	tmdb *tmdb.Service
	tvdb *tvdb.Service

	files map[string]os.FileInfo
}

func New(args Args, tmdb *tmdb.Service, tvdb *tvdb.Service) *Namer {
	return &Namer{
		args:  args,
		tmdb:  tmdb,
		tvdb:  tvdb,
		files: map[string]os.FileInfo{},
	}
}

func (pn *Namer) Run() error {
	if err := pn.collectFiles(); err != nil {
		return err
	}
	return nil
}

func (pn *Namer) collectFiles() error {
	if err := filepath.Walk(pn.args.Path, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			pn.files[path] = node
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}
