package main

import (
	"flag"
	"github.com/florianehmke/plexname/parser"
	"log"
	"strings"

	"github.com/florianehmke/plexname"
	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	pn := plexname.New(
		parseArgs(),
		tmdb.NewService(tmdb.BaseURL, config.GetToken("tmdb")),
		tvdb.NewService(tvdb.BaseURL, config.GetToken("tvdb")),
	)
	pn.PrintPlexName()
}

func parseArgs() plexname.Args {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("no release name given")
	}
	query := strings.Join(flag.Args(), " ")

	return plexname.Args{
		Query:     query,
		Overrides: parser.Result{},
	}
}
