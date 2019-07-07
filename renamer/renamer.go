package renamer

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

type Renamer struct {
	args Args

	searcher search.Searcher
	fs       fs.FileSystem

	files []fileInfo
}

func New(args Args, searcher search.Searcher, fs fs.FileSystem) *Renamer {
	return &Renamer{
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

func (r *Renamer) Run() error {
	if info, err := os.Stat(r.args.SourcePath); err == nil {
		if info.IsDir() {
			return r.runDir()
		} else {
			return r.runFile()
		}
	} else {
		return err
	}
}

func (r *Renamer) runDir() error {
	if err := r.collectFiles(); err != nil {
		return err
	}
	if err := r.collectNewPaths(); err != nil {
		return err
	}
	if err := r.moveAndRename(); err != nil {
		return err
	}
	return nil
}

func (r *Renamer) collectFiles() error {
	if err := filepath.Walk(r.args.SourcePath, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			p := filepath.ToSlash(path)
			r.files = append(r.files, fileInfo{
				currentFilePath:         p,
				currentRelativeFilePath: strings.TrimPrefix(p, r.args.SourcePath+"/"),
			})
		}
		return nil
	}); err != nil {
		return fmt.Errorf("directory scan failed: %v", err)
	}
	return nil
}

func (r *Renamer) collectNewPaths() error {
	var files []fileInfo
	for _, f := range r.files {
		log.Info(fmt.Sprintf("Processing: %s", f.currentFilePath))
		pr := parser.Parse(f.currentFilePath, f.currentFilePath, r.args.OnlyFile, r.args.OnlyDir, r.args.Overrides)
		sr, err := r.search(pr)

		if err != nil {
			return fmt.Errorf("search for %s failed: %v", f.currentRelativeFilePath, err)
		}

		plexName, err := plexName(pr, &sr)
		if err != nil {
			return fmt.Errorf("could not get a plex name for %s: %v", f.currentRelativeFilePath, err)
		}

		if pr.IsMovie() {
			// The new directory..
			f.newPath = r.args.TargetPath + "/" + plexName

			// .. and the filename inside of that directory.
			// See: https://support.plex.tv/articles/200381043-multi-version-movies/
			extension := strings.ToLower(filepath.Ext(f.fileName()))
			versionInfo := pr.VersionInfo()
			fileName := joinNonEmpty(" - ", plexName, versionInfo)
			f.newFilePath = f.newPath + "/" + fileName + extension
		}

		if pr.IsTV() {
			// The new directory + Season Folder ...
			f.newPath = fmt.Sprintf("%s/%s/Season %02d", r.args.TargetPath, plexName, pr.Season)

			// .. and the episode filename inside of that directory.
			// See: https://support.plex.tv/articles/naming-and-organizing-your-tv-show-files/
			extension := strings.ToLower(filepath.Ext(f.fileName()))
			versionInfo := pr.VersionInfo()
			tvInfo := joinNonEmpty("", toSeasonString(pr.Season), toEpisodeString(pr.Episode1), toEpisodeString(pr.Episode2))
			fileName := joinNonEmpty(" - ", plexName, tvInfo, versionInfo)
			f.newFilePath = f.newPath + "/" + fileName + extension
		}
		files = append(files, f)
	}
	r.files = files
	return nil
}

func (r *Renamer) moveAndRename() error {
	for _, f := range r.files {
		if err := r.move(f.currentFilePath, f.newFilePath); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renamer) runFile() error {
	log.Info(fmt.Sprintf("Processing: %s", r.args.SourcePath))
	dir, file := filepath.Split(r.args.SourcePath)
	pr := parser.Parse(file, file, r.args.OnlyFile, r.args.OnlyDir, r.args.Overrides)
	sr, err := r.search(pr)
	if err != nil {
		return fmt.Errorf("search for %s failed: %v", file, err)
	}

	plexName, err := plexName(pr, &sr)
	if err != nil {
		return fmt.Errorf("could not get a plex name for %s: %v", r.args.SourcePath, err)
	}

	var newFilePath string
	if pr.IsMovie() {
		extension := strings.ToLower(filepath.Ext(file))
		versionInfo := pr.VersionInfo()
		fileName := joinNonEmpty(" - ", plexName, versionInfo)
		newFilePath = dir + fileName + extension
	}

	if pr.IsTV() {
		extension := strings.ToLower(filepath.Ext(file))
		versionInfo := pr.VersionInfo()
		tvInfo := joinNonEmpty("", toSeasonString(pr.Season), toEpisodeString(pr.Episode1), toEpisodeString(pr.Episode2))
		fileName := joinNonEmpty(" - ", plexName, tvInfo, versionInfo)
		newFilePath = dir + fileName + extension
	}

	return r.move(r.args.SourcePath, newFilePath)
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

func (r *Renamer) search(pr *parser.Result) (search.Result, error) {
	if pr.IsMovie() {
		return r.searcher.SearchMovie(search.Query{Title: pr.Title, Year: pr.Year})
	}
	if pr.IsTV() {
		return r.searcher.SearchTV(search.Query{Title: pr.Title, Year: pr.Year})
	}
	return search.Result{}, errors.New("can not search for unknown media type")
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.currentRelativeFilePath)
	return fileName
}

func (r *Renamer) move(source, target string) error {
	if r.skipBasedOnExtension(source) {
		log.Warn(fmt.Sprintf("Skipping %s (based on extension)", source))
		return nil
	}

	newDir, fileName := filepath.Split(target)

	osNewDir := filepath.FromSlash(newDir)
	if err := r.fs.MkdirAll(osNewDir); err != nil {
		return fmt.Errorf("mkdir of %s failed: %v", osNewDir, err)
	}

	osTarget := filepath.FromSlash(target)
	osSource := filepath.FromSlash(source)
	if err := r.fs.Rename(osSource, osTarget); err != nil {
		return fmt.Errorf("move of %s to %s failed: %v", fileName, osNewDir, err)
	}

	log.Info(fmt.Sprintf("Renamed:\nSource: %s\nTarget: %s", source, target))
	return nil
}

func (r *Renamer) skipBasedOnExtension(s string) bool {
	if len(r.args.Extensions) == 0 {
		return false
	}
	skip := true
	ext := strings.TrimLeft(filepath.Ext(s), ".")
	for _, e := range r.args.Extensions {
		if e == ext {
			skip = false
		}
	}
	return skip
}

func toEpisodeString(ep int) string {
	if ep != 0 {
		return fmt.Sprintf("E%02d", ep)
	}
	return ""
}

func toSeasonString(s int) string {
	if s != 0 {
		return fmt.Sprintf("S%02d", s)
	}
	return ""
}

func joinNonEmpty(sep string, slist ...string) string {
	var slice []string
	for _, s := range slist {
		if s != "" {
			slice = append(slice, s)
		}
	}
	return strings.Join(slice, sep)
}
