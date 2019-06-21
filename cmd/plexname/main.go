package main

import (
	"flag"
	"log"
	"strings"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	pn := namer.New(
		parseArgs(),
		tmdb.NewService(tmdb.BaseURL, config.GetToken("tmdb")),
		tvdb.NewService(tvdb.BaseURL, config.GetToken("tvdb")),
	)
	pn.PrintPlexName()
}

func parseArgs() namer.Args {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no release name given")
	}
	query := strings.Join(flag.Args(), " ")

	return namer.Args{
		Query:     query,
		Overrides: parser.Result{},
	}
}
