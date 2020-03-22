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
	params Parameters

	searcher search.Searcher
	fs       fs.FileSystem

	files []fileInfo
}

func New(args Parameters, searcher search.Searcher, fs fs.FileSystem) *Renamer {
	return &Renamer{
		params:   args,
		searcher: searcher,
		fs:       fs,
		files:    []fileInfo{},
	}
}

type fileInfo struct {
	currentFilePath string

	newPath     string
	newFilePath string
}

func (r *Renamer) parse(source, target string) parser.Result {
	srcPath, srcFile := filepath.Split(source)
	_, srcDir := filepath.Split(strings.TrimRight(srcPath, "/"))
	srcDirAndFile := srcDir + "/" + srcFile

	// TODO guess media type

	var toParse string
	if r.params.OnlyFile {
		toParse = srcFile
	} else if r.params.OnlyDir {
		toParse = srcDir
	} else {
		toParse = srcDirAndFile
	}

	return parser.Parse(toParse, r.params.Overrides)
}

func (r *Renamer) Run() error {
	if info, err := os.Stat(r.params.SourcePath); err == nil {
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
	if err := filepath.Walk(r.params.SourcePath, func(path string, node os.FileInfo, err error) error {
		if !node.IsDir() {
			p := filepath.ToSlash(path)
			r.files = append(r.files, fileInfo{currentFilePath: p})
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
		pr := r.parse(f.currentFilePath, r.params.TargetPath)

		sr, err := r.search(pr)
		if err != nil {
			return fmt.Errorf("search for %s failed: %v", f.currentFilePath, err)
		}

		plexName, err := plexName(pr, sr)
		if err != nil {
			return fmt.Errorf("could not get a plex name for %s: %v", f.currentFilePath, err)
		}

		newPath, err := newDirectoryPath(r.params.TargetPath, plexName, pr)
		if err != nil {
			return fmt.Errorf("could not create directory path for %s: %v", f.currentFilePath, err)
		}
		f.newPath = newPath

		newFilePath, err := newFilePath(newPath, f.fileName(), plexName, pr)
		if err != nil {
			return fmt.Errorf("could not create file path for %s: %v", f.currentFilePath, err)
		}
		f.newFilePath = newFilePath

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
	log.Info(fmt.Sprintf("Processing: %s", r.params.SourcePath))
	dir, file := filepath.Split(r.params.SourcePath)

	pr := r.parse(file, file)

	sr, err := r.search(pr)
	if err != nil {
		return fmt.Errorf("search for %s failed: %v", file, err)
	}

	plexName, err := plexName(pr, sr)
	if err != nil {
		return fmt.Errorf("could not get a plex name for %s: %v", r.params.SourcePath, err)
	}

	newFilePath, err := newFilePath(dir, file, plexName, pr)
	if err != nil {
		return fmt.Errorf("could not create file path for %s: %v", r.params.SourcePath, err)
	}

	return r.move(r.params.SourcePath, newFilePath)
}

func plexName(pr parser.Result, sr search.Result) (string, error) {
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

func newFilePath(base string, oldFileName string, plexName string, pr parser.Result) (string, error) {
	base = strings.TrimRight(base, "/")
	extension := strings.ToLower(filepath.Ext(oldFileName))
	versionInfo := versionInfo(pr)
	if pr.IsTV() {
		tvInfo := tvInfo(pr)
		fileName := joinNonEmpty(" - ", plexName, tvInfo, versionInfo)
		return base + "/" + fileName + extension, nil
	}
	if pr.IsMovie() {
		fileName := joinNonEmpty(" - ", plexName, versionInfo)
		return base + "/" + fileName + extension, nil
	}
	return "", errors.New("can't create file path for unknown media type")
}

func newDirectoryPath(base string, plexName string, pr parser.Result) (string, error) {
	base = strings.TrimRight(base, "/")
	if pr.IsTV() {
		return fmt.Sprintf("%s/%s/Season %02d", base, plexName, pr.Season), nil
	}
	if pr.IsMovie() {
		return base + "/" + plexName, nil
	}
	return "", errors.New("can't create directory path for unknown media type")
}

func (r *Renamer) search(pr parser.Result) (search.Result, error) {
	if pr.IsMovie() {
		return r.searcher.SearchMovie(search.Query{Title: pr.Title, Year: pr.Year})
	}
	if pr.IsTV() {
		return r.searcher.SearchTV(search.Query{Title: pr.Title, Year: pr.Year})
	}
	return search.Result{}, errors.New("can not search for unknown media type")
}

func (fi *fileInfo) fileName() string {
	_, fileName := path.Split(fi.currentFilePath)
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

	log.Info(fmt.Sprintf("Renamed:\nSource: %s\nTarget: %s", osSource, osTarget))
	return nil
}

func (r *Renamer) skipBasedOnExtension(s string) bool {
	if len(r.params.Extensions) == 0 {
		return false
	}
	skip := true
	ext := strings.TrimLeft(filepath.Ext(s), ".")
	for _, e := range r.params.Extensions {
		if e == ext {
			skip = false
		}
	}
	return skip
}

func toEpisodeString(ep int) string {
	if ep > 0 {
		return fmt.Sprintf("E%02d", ep)
	}
	return ""
}

func toSeasonString(s int, special parser.ParseBool) string {
	if s > 0 || special == parser.True {
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

func tvInfo(pr parser.Result) string {
	return joinNonEmpty("",
		toSeasonString(pr.Season, pr.Special),
		toEpisodeString(pr.Episode1),
		toEpisodeString(pr.Episode2),
	)
}

func versionInfo(pr parser.Result) string {
	tokens := []string{}
	if pr.Language != parser.LangNA {
		tokens = append(tokens, pr.Language.String())
	}
	if pr.Resolution != parser.ResNA {
		tokens = append(tokens, pr.Resolution.String())
	}
	if pr.DualLanguage == parser.True {
		tokens = append(tokens, "DL")
	}
	if pr.Source != parser.SourceNA {
		tokens = append(tokens, pr.Source.String())
	}
	if pr.Remux == parser.True {
		tokens = append(tokens, "Remux")
	}
	return strings.Join(tokens, ".")
}
