package namer

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/prompt"
	"github.com/florianehmke/plexname/search"
)

type Args struct {
	Path      string
	Overrides parser.Result
}

type Namer struct {
	args Args

	searcher search.Searcher
	prompter prompt.Prompter
	fs       fs.FileSystem

	files []fileInfo
}

func New(args Args, searcher search.Searcher, prompter prompt.Prompter, fs fs.FileSystem) *Namer {
	args.Path = convertPathToSlash(args.Path)
	return &Namer{
		args:     args,
		searcher: searcher,
		prompter: prompter,
		fs:       fs,
		files:    []fileInfo{},
	}
}

type fileInfo struct {
	currentFilePath         string
	currentRelativeFilePath string

	newPath     string
	newFilePath string
}

func (n *Namer) Run() error {
	if err := n.collectFiles(); err != nil {
		return err
	}
	if err := n.collectNewPaths(); err != nil {
		return err
	}
	if err := n.moveAndRename(); err != nil {
		return err
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

func (fi *fileInfo) segmentToParse() string {
	segments := strings.Split(fi.currentRelativeFilePath, "/")
	segment, length := "", 0
	for _, s := range segments {
		if len(s) > length {
			segment, length = s, len(s)
		}
	}
	return segment
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.currentRelativeFilePath)
	return fileName
}

func (n *Namer) collectFiles() error {
	if err := filepath.Walk(n.args.Path, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			p := filepath.ToSlash(path)
			n.files = append(n.files, fileInfo{
				currentFilePath:         p,
				currentRelativeFilePath: strings.TrimPrefix(p, n.args.Path+"/"),
			})
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}

func (n *Namer) collectNewPaths() error {
	for i, _ := range n.files {
		f := &n.files[i]
		pr := parser.Parse(f.segmentToParse(), n.args.Overrides)
		sr, err := n.Search(pr)

		if err != nil {
			return fmt.Errorf("search for %s failed: %v", f.currentRelativeFilePath, err)
		}
		if len(sr) == 0 {
			return fmt.Errorf("no search result title %s of %s", pr.Title, f.currentRelativeFilePath)
		}

		var result *search.Result
		if len(sr) > 1 {
			choices := []string{"Multiple results found online, pick one of:"}
			for i, r := range sr {
				choices = append(choices, fmt.Sprintf("[%d] %s (%d)", i+1, r.Title, r.Year))
			}
			i, err := n.prompter.AskNumber(strings.Join(choices, "\n"))
			if err != nil {
				return fmt.Errorf("prompt error: %v", err)
			}
			result = &sr[i-1]
		} else {
			result = &sr[0]
		}

		newName, err := plexName(pr, result)
		if err != nil {
			return fmt.Errorf("could not get a plex name for %s: %v", f.currentRelativeFilePath, err)
		}

		f.newPath = n.args.Path + "/" + newName
		f.newFilePath = f.newPath + "/" + f.fileName()
	}
	return nil
}

func (n *Namer) moveAndRename() error {
	for _, f := range n.files {
		osNewPath := filepath.FromSlash(f.newPath)
		if err := n.fs.MkdirAll(osNewPath); err != nil {
			return fmt.Errorf("mkdir of %s failed: %v", osNewPath, err)
		}

		osNewFilePath := filepath.FromSlash(f.newFilePath)
		osOldFilePath := filepath.FromSlash(f.currentFilePath)
		if err := n.fs.Rename(osOldFilePath, osNewFilePath); err != nil {
			return fmt.Errorf("move of %s to %s failed: %v", f.fileName(), osNewFilePath, err)
		}
	}
	return nil
}

func convertPathToSlash(path string) string {
	slashPath := filepath.ToSlash(path)
	slashPath = strings.TrimRight(slashPath, "/")
	return slashPath
}
