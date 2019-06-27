package main

import (
	"flag"
	"log"
	"strings"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	pn := namer.New(
		parseArgs(),
		tmdb.NewClient(tmdb.BaseURL, config.GetToken("tmdb")),
		tvdb.NewClient(tvdb.BaseURL, config.GetToken("tvdb")),
		fs.NewFileSystem(),
	)
	pn.Run()
}

func parseArgs() namer.Args {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no directory given")
	}
	path := strings.Join(flag.Args(), " ")

	return namer.Args{
		Path:      path,
		Overrides: parser.Result{},
	}
}
