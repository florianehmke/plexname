package renamer

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/florianehmke/plexname/parser"
)

type Args struct {
	sourcePath string
	targetPath string
	overrides  parser.Result
	extensions []string
	dryRun     bool
}

func NewArgs(source, target string, overrides parser.Result, extensions []string, dryRun bool) Args {
	args := Args{}
	args.sourcePath = filepath.ToSlash(source)
	if target == "" {
		args.targetPath = args.sourcePath
	} else {
		args.targetPath = filepath.ToSlash(target)
	}
	args.sourcePath = strings.TrimRight(args.sourcePath, "/")
	args.targetPath = strings.TrimRight(args.targetPath, "/")
	args.overrides = overrides
	args.extensions = extensions
	args.dryRun = dryRun
	return args
}

func GetArgsFromFlags() Args {
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
	if flag.NArg() == 1 {
		return NewArgs(flag.Arg(0), flag.Arg(0), overrides, extSliceFor(extensions), dryRun)
	} else {
		return NewArgs(flag.Arg(0), flag.Arg(1), overrides, extSliceFor(extensions), dryRun)
	}
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

func extSliceFor(s string) []string {
	var slice []string
	if s != "" {
		slice = strings.Split(s, ",")
	}
	return slice
}

func (a *Args) SourcePath() string {
	return a.sourcePath
}

func (a *Args) TargetPath() string {
	return a.targetPath
}

func (a *Args) Overrides() parser.Result {
	return a.overrides
}

func (a *Args) Extensions() []string {
	return a.extensions
}

func (a *Args) DryRun() bool {
	return a.dryRun
}
