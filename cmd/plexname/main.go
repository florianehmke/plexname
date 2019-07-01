package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/log"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/prompt"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	pn := namer.New(
		parseArgs(),
		search.NewSearcher(
			tmdb.NewClient(tmdb.BaseURL, config.GetToken("tmdb")),
			tvdb.NewClient(tvdb.BaseURL, config.GetToken("tvdb")),
		),
		prompt.NewPrompter(),
		fs.NewFileSystem(),
	)
	parseArgs()
	if err := pn.Run(); err != nil {
		log.Error(fmt.Sprintf("rename failed: %v", err))
		os.Exit(2)
	}
	log.Info("Yay, done!")
	os.Exit(0)
}

func parseArgs() namer.Args {
	flag.Usage = usage
	overrides := parser.Result{}
	flag.StringVar(&overrides.Title, "title", "", "movie/tv title")
	flag.IntVar(&overrides.Year, "year", 0, "movie/tv year of release")
	flag.IntVar(&overrides.Season, "season", 0, "tv season of release")
	flag.IntVar(&overrides.Year, "episode", 0, "tv episode of release")

	// TODO:
	// MediaType    MediaType
	// Resolution   Resolution
	// Source       Source
	// Language     Language
	// Remux        bool
	// Proper       bool
	// DualLanguage bool

	flag.Parse()

	if flag.NArg() == 0 || flag.NArg() > 2 {
		flag.Usage()
		os.Exit(1)
	}
	if flag.NArg() == 1 {
		return namer.NewArgs(flag.Arg(0), flag.Arg(0), overrides)
	} else {
		return namer.NewArgs(flag.Arg(0), flag.Arg(1), overrides)
	}
}

func usage() {
	fmt.Println("plexname")
	fmt.Println("  Rename your media files and folders for the Plex Media Server.")
	fmt.Println()
	fmt.Println("Usage: ")
	fmt.Println("  plexname [option]... source [target]")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
