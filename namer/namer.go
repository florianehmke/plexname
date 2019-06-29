package namer

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/log"
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

func (n *Namer) Run() error {
	if err := n.collectFiles(); err != nil {
		return err
	}

	for p, f := range n.files {
		pr := parser.Parse(f.segmentToParse(), n.args.Overrides)
		sr, err := n.Search(pr)
		if err != nil {
			return fmt.Errorf("search for %s failed: %v", f.relativePath, err)
		}

		if len(sr) == 0 {
			return fmt.Errorf("no search result title %s of %s", pr.Title, f.relativePath)
		}
		if len(sr) > 1 {
			log.Warn(fmt.Sprintf("ambigious result for %s", pr.Title))
		}

		newName, err := plexName(pr, &sr[0])
		if err != nil {
			return fmt.Errorf("could not get a plex name for %s: %v", f.relativePath, err)
		}

		newPath := n.args.Path + string(os.PathSeparator) + newName
		if err := n.fs.MkdirAll(newPath); err != nil {
			return fmt.Errorf("mkdir of %s failed: %v", newPath, err)
		}
		newFilePath := newPath + string(os.PathSeparator) + f.fileName()
		if err := n.fs.Rename(p, newFilePath); err != nil {
			return fmt.Errorf("move of %s to %s failed: %v", f.fileName(), newFilePath, err)
		}
	}
	return nil
}

func plexName(pr *parser.Result, sr *search.Result) (string, error) {
	year := pr.Year
	if year == 0 {
		year = sr.Year
	}
	if year == 0 {
		return "", errors.New("neither parser nor search yielded a year")
	}
	return fmt.Sprintf("%s (%d)", sr.Title, year), nil
}

func (n *Namer) Search(pr *parser.Result) ([]search.Result, error) {
	if pr.IsMovie() {
		return n.searcher.SearchMovie(pr.Title, pr.Year)
	}
	if pr.IsTV() {
		return n.searcher.SearchTV(pr.Title, pr.Year)
	}
	return nil, errors.New("can not search for unknown media type")
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
			segment, length = s, len(s)
		}
	}
	return segment
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.relativePath)
	return fileName
}

func (n *Namer) collectFiles() error {
	if err := filepath.Walk(n.args.Path, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			n.files[path] = fileInfo{
				relativePath: strings.TrimPrefix(path, n.args.Path+string(os.PathSeparator)),
				info:         node,
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}
