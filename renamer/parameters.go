package renamer

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/parser"
)

type Parameters struct {
	SourcePath string
	TargetPath string
	Overrides  parser.Result
	Extensions []string

	DryRun bool

	OnlyFile bool
	OnlyDir  bool
}

func NewParameters(source, target string, overrides parser.Result, extensions []string, dryRun, onlyFile, onlyDir bool) Parameters {
	targetPath := target
	if targetPath == "" {
		targetPath = source
	}

	return Parameters{
		SourcePath: strings.TrimRight(filepath.ToSlash(source), "/"),
		TargetPath: strings.TrimRight(filepath.ToSlash(targetPath), "/"),
		Overrides:  overrides,
		Extensions: extensions,
		DryRun:     dryRun,
		OnlyFile:   onlyFile,
		OnlyDir:    onlyDir,
	}
}

func GetParametersFromFlags() Parameters {
	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "do a dry run")

	overrides := parser.Result{}
	flag.StringVar(&overrides.Title, "title", "", "movie/tv title")
	flag.IntVar(&overrides.Year, "year", 0, "movie/tv year of release")
	flag.IntVar(&overrides.Season, "season", 0, "tv season of release")
	flag.IntVar(&overrides.Episode1, "episode1", 0, "tv episode1 of release")
	flag.IntVar(&overrides.Episode2, "episode2", 0, "tv episode2 of release")
	var proper, remux, dualLang string
	flag.StringVar(&proper, "proper", "", "proper release")
	flag.StringVar(&remux, "remux", "", "remux of source, no encode")
	flag.StringVar(&dualLang, "dl", "", "dual language")

	var mediaType, resolution, source, lang string
	flag.StringVar(&mediaType, "media-type", "", "media type (movie|tv)")
	flag.StringVar(&resolution, "resolution", "", "resolution: 720p, 1080p etc")
	flag.StringVar(&source, "source", "", "media source (web-dl, blu-ray etc)")
	flag.StringVar(&lang, "lang", "", "audio language")

	var extensions string
	flag.StringVar(&extensions, "extensions", "", "move only file with the given extension")

	var onlyDir, onlyFile bool
	flag.BoolVar(&onlyDir, "only-dir", false, "parse only the directory name")
	flag.BoolVar(&onlyFile, "only-file", false, "parse only file name")
	flag.Parse()

	overrides.Proper = boolFor(proper)
	overrides.Remux = boolFor(remux)
	overrides.DualLanguage = boolFor(dualLang)

	overrides.MediaType = mediaTypeFor(mediaType)
	overrides.Resolution = resolutionFor(resolution)
	overrides.Source = sourceFor(source)
	overrides.Language = languageFor(lang)

	if flag.NArg() == 0 || flag.NArg() > 2 {
		flag.Usage()
		os.Exit(1)
	}

	var sourcePath, targetPath string
	var err error

	sourcePath = flag.Arg(0)
	sourcePath, err = filepath.Abs(sourcePath)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		targetPath, err = os.Getwd()
		if err != nil {
			flag.Usage()
			os.Exit(1)
		}
	} else {
		targetPath = flag.Arg(1)
	}

	return NewParameters(sourcePath, targetPath, overrides, splitExtensions(extensions), dryRun, onlyFile, onlyDir)

}

func mediaTypeFor(s string) parser.MediaType {
	mt, err := parser.ParseMediaType(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return mt
}

func resolutionFor(s string) parser.Resolution {
	r, err := parser.ParseResolution(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return r
}

func sourceFor(s string) parser.Source {
	src, err := parser.ParseSource(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return src
}

func languageFor(s string) parser.Language {
	l, err := parser.ParseLanguage(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return l
}

func boolFor(s string) parser.ParseBool {
	ls := strings.ToLower(s)
	if ls == "true" {
		return parser.True
	}
	if ls == "false" {
		return parser.False
	}
	return parser.Unknown
}

func splitExtensions(s string) []string {
	var slice []string
	if s != "" {
		slice = strings.Split(s, ",")
	}
	return slice
}
