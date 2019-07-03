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
	SourcePath string
	TargetPath string
	Overrides  parser.Result
}

type Namer struct {
	args Args

	searcher search.Searcher
	fs       fs.FileSystem

	files []fileInfo
}

func NewArgs(source, target string, overrides parser.Result) Args {
	args := Args{}
	args.SourcePath = filepath.ToSlash(source)
	if target == "" {
		args.TargetPath = args.SourcePath
	} else {
		args.TargetPath = filepath.ToSlash(target)
	}
	args.SourcePath = strings.TrimRight(args.SourcePath, "/")
	args.TargetPath = strings.TrimRight(args.TargetPath, "/")
	args.Overrides = overrides
	return args
}

func New(args Args, searcher search.Searcher, fs fs.FileSystem) *Namer {
	return &Namer{
		args:     args,
		searcher: searcher,
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
	if info, err := os.Stat(n.args.SourcePath); err == nil {
		if info.IsDir() {
			return n.runDir()
		} else {
			return n.runFile()
		}
	} else {
		return err
	}
}

func (n *Namer) runDir() error {
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

func (n *Namer) runFile() error {
	log.Info(fmt.Sprintf("Processing: %s", n.args.SourcePath))
	dir, file := filepath.Split(n.args.SourcePath)
	pr := parser.Parse(file, file, n.args.Overrides)
	sr, err := n.search(pr)
	if err != nil {
		return fmt.Errorf("search for %s failed: %v", file, err)
	}

	plexName, err := plexName(pr, &sr)
	if err != nil {
		return fmt.Errorf("could not get a plex name for %s: %v", n.args.SourcePath, err)
	}

	newFilePath := ""
	if pr.IsMovie() {
		// .. and the filename inside of that directory.
		// See: https://support.plex.tv/articles/200381043-multi-version-movies/
		extension := strings.ToLower(filepath.Ext(file))
		versionInfo := pr.VersionInfo()
		fileName := fmt.Sprintf("%s - %s%s", plexName, versionInfo, extension)
		newFilePath = dir + fileName
	}

	if pr.IsTV() {
		// .. and the episode filename inside of that directory.
		// See: https://support.plex.tv/articles/naming-and-organizing-your-tv-show-files/
		extension := strings.ToLower(filepath.Ext(file))
		versionInfo := pr.VersionInfo()
		fileName := fmt.Sprintf("%s - s%02de%02d - %s%s", plexName, pr.Season, pr.Episode, versionInfo, extension)
		newFilePath = dir + fileName
	}

	return n.move(n.args.SourcePath, newFilePath)
}

func plexName(pr *parser.Result, sr *search.Result) (string, error) {
	year := pr.Year
	if year == 0 {
		year = sr.Year
	}
	if year == 0 {
		if pr.IsMovie() {
			return "", errors.New("neither parser nor search yielded a year")
		}
		if pr.IsTV() {
			// For TV it is okay if year is missing.
			return sr.Title, nil
		}
	}
	return fmt.Sprintf("%s (%d)", sr.Title, year), nil
}

func (n *Namer) search(pr *parser.Result) (search.Result, error) {
	if pr.IsMovie() {
		return n.searcher.SearchMovie(search.Query{Title: pr.Title, Year: pr.Year})
	}
	if pr.IsTV() {
		return n.searcher.SearchTV(search.Query{Title: pr.Title, Year: pr.Year})
	}
	return search.Result{}, errors.New("can not search for unknown media type")
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.currentRelativeFilePath)
	return fileName
}

func (n *Namer) collectFiles() error {
	if err := filepath.Walk(n.args.SourcePath, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			p := filepath.ToSlash(path)
			n.files = append(n.files, fileInfo{
				currentFilePath:         p,
				currentRelativeFilePath: strings.TrimPrefix(p, n.args.SourcePath+"/"),
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
		log.Info(fmt.Sprintf("Processing: %s", f.currentFilePath))
		pr := parser.Parse(f.currentFilePath, f.currentFilePath, n.args.Overrides)
		sr, err := n.search(pr)

		if err != nil {
			return fmt.Errorf("search for %s failed: %v", f.currentRelativeFilePath, err)
		}

		plexName, err := plexName(pr, &sr)
		if err != nil {
			return fmt.Errorf("could not get a plex name for %s: %v", f.currentRelativeFilePath, err)
		}

		if pr.IsMovie() {
			// The new directory..
			f.newPath = n.args.TargetPath + "/" + plexName

			// .. and the filename inside of that directory.
			// See: https://support.plex.tv/articles/200381043-multi-version-movies/
			extension := strings.ToLower(filepath.Ext(f.fileName()))
			versionInfo := pr.VersionInfo()
			fileName := fmt.Sprintf("%s - %s%s", plexName, versionInfo, extension)
			f.newFilePath = f.newPath + "/" + fileName
		}

		if pr.IsTV() {
			// The new directory + Season Folder ...
			f.newPath = fmt.Sprintf("%s/%s/Season %02d", n.args.TargetPath, plexName, pr.Season)

			// .. and the episode filename inside of that directory.
			// See: https://support.plex.tv/articles/naming-and-organizing-your-tv-show-files/
			extension := strings.ToLower(filepath.Ext(f.fileName()))
			versionInfo := pr.VersionInfo()
			fileName := fmt.Sprintf("%s - s%02de%02d - %s%s", plexName, pr.Season, pr.Episode, versionInfo, extension)
			f.newFilePath = f.newPath + "/" + fileName
		}
	}
	return nil
}

func (n *Namer) moveAndRename() error {
	for _, f := range n.files {
		if err := n.move(f.currentFilePath, f.newFilePath); err != nil {
			return err
		}
	}
	return nil
}

func (n *Namer) move(source, target string) error {
	newDir, fileName := filepath.Split(target)

	osNewDir := filepath.FromSlash(newDir)
	if err := n.fs.MkdirAll(osNewDir); err != nil {
		return fmt.Errorf("mkdir of %s failed: %v", osNewDir, err)
	}

	osTarget := filepath.FromSlash(target)
	osSource := filepath.FromSlash(source)
	if err := n.fs.Rename(osSource, osTarget); err != nil {
		return fmt.Errorf("move of %s to %s failed: %v", fileName, osNewDir, err)
	}

	log.Info(fmt.Sprintf("Renamed to: %s", target))
	return nil
}
