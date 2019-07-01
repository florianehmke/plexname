package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	if err := pn.Run(); err != nil {
		log.Error(fmt.Sprintf("rename failed: %v", err))
		os.Exit(1)
	}
	log.Info("Yay, done!")
	os.Exit(0)
}

func parseArgs() namer.Args {
	overrides := parser.Result{}
	flag.StringVar(&overrides.Title, "title", "", "movie/tv title")
	flag.IntVar(&overrides.Year, "year", 0, "movie/tv year of release")
	flag.IntVar(&overrides.Season, "season", 0, "tv season of release")
	flag.IntVar(&overrides.Year, "episode", 0, "tv episode of release")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no directory given")
	}
	path := strings.Join(flag.Args(), " ")

	return namer.NewArgs(path, path, overrides)
}
