package namer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/search"
)

type Args struct {
	Path      string
	Overrides parser.Result
}

type Namer struct {
	args Args

	searcher search.Searcher
	fs       fs.FileSystem

	files map[string]fileInfo
}

func New(args Args, searcher search.Searcher, fs fs.FileSystem) *Namer {
	return &Namer{
		args:     args,
		searcher: searcher,
		fs:       fs,
		files:    map[string]fileInfo{},
	}
}

func (pn *Namer) Run() error {
	if err := pn.collectFiles(); err != nil {
		return err
	}

	// FIXME needs refactor
	for p, f := range pn.files {
		pr := parser.Parse(f.segmentToParse(), pn.args.Overrides)

		newName := ""
		if pr.IsMovie() {
			result, _ := pn.searcher.SearchMovie(pr.Title, pr.Year)
			newName = result[0].Title
		}
		if pr.IsTV() {
			result, _ := pn.searcher.SearchTV(pr.Title, pr.Year)
			newName = result[0].Title
		}

		if newName != "" {
			if err := pn.fs.MkdirAll(newName); err != nil {
				return err
			}
			if err := pn.fs.Rename(p, pn.args.Path+string(os.PathSeparator)+newName+string(os.PathSeparator)+f.fileName()); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no suitable name found for: %s", f.relativePath)
		}
	}
	return nil
}

type fileInfo struct {
	relativePath string
	info         os.FileInfo
}

func (fi *fileInfo) segmentToParse() string {
	segments := strings.Split(fi.relativePath, string(os.PathSeparator))
	segment, length := "", 0
	for _, s := range segments {
		if len(s) > length {
			segment, length = s, len(segment)
		}
	}
	return segment
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.relativePath)
	return fileName
}

func (pn *Namer) collectFiles() error {
	if err := filepath.Walk(pn.args.Path, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			pn.files[path] = fileInfo{
				relativePath: strings.TrimPrefix(path, pn.args.Path+string(os.PathSeparator)),
				info:         node,
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}
