package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/log"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/prompt"
	"github.com/florianehmke/plexname/renamer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	args, dryRun := parseArgs()
	rn := renamer.New(
		args,
		search.NewSearcher(
			tmdb.NewClient(tmdb.BaseURL, config.GetToken("tmdb")),
			tvdb.NewClient(tvdb.BaseURL, config.GetToken("tvdb")),
			prompt.NewPrompter(),
		),
		fs.NewFileSystem(dryRun),
	)
	if err := rn.Run(); err != nil {
		log.Error(fmt.Sprintf("rename failed: %v", err))
		os.Exit(1)
	}
	log.Info("Yay, done!")
	os.Exit(0)
}

func parseArgs() (renamer.Args, bool) {
	flag.Usage = usage

	var dryRun bool
	flag.BoolVar(&dryRun, "dry", false, "do a dry run")

	overrides := parser.Result{}
	flag.StringVar(&overrides.Title, "title", "", "movie/tv title")
	flag.IntVar(&overrides.Year, "year", 0, "movie/tv year of release")
	flag.IntVar(&overrides.Season, "season", 0, "tv season of release")
	flag.IntVar(&overrides.Year, "episode", 0, "tv episode of release")
	var proper, remux, dualLang string
	flag.StringVar(&proper, "proper", "", "proper release")
	flag.StringVar(&remux, "remux", "", "remux of source, no encode")
	flag.StringVar(&dualLang, "dl", "", "dual language")
	overrides.Proper = boolFor(proper)
	overrides.Remux = boolFor(remux)
	overrides.DualLanguage = boolFor(dualLang)

	var mediaType, resolution, source, lang string
	flag.StringVar(&mediaType, "media-type", "", "media type (movie|tv)")
	flag.StringVar(&resolution, "resolution", "", "resolution: 720p, 1080p etc")
	flag.StringVar(&source, "source", "", "media source (web-dl, blu-ray etc)")
	flag.StringVar(&lang, "lang", "", "audio language")

	var extensions string
	flag.StringVar(&extensions, "extensions", "", "move only file with the given extension")
	flag.Parse()

	overrides.MediaType = mediaTypeFor(mediaType)
	overrides.Resolution = resolutionFor(resolution)
	overrides.Source = sourceFor(source)
	overrides.Language = languageFor(lang)

	if flag.NArg() == 0 || flag.NArg() > 2 {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NArg() == 1 {
		return renamer.NewArgs(flag.Arg(0), flag.Arg(0), overrides, extSliceFor(extensions)), dryRun
	} else {
		return renamer.NewArgs(flag.Arg(0), flag.Arg(1), overrides, extSliceFor(extensions)), dryRun
	}
}

func usage() {
	fmt.Println("plexname")
	fmt.Println("  Rename your media files and folders for the Plex Media Server.")
	fmt.Println()
	fmt.Println("Usage: ")
	fmt.Println("  plexname [option]... source-dir [target-dir]")
	fmt.Println("  plexname [option]... file")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  plexname -extensions=mkv,mp4 -lang english -remux downloads movies")
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
