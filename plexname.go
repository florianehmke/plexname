package plexname

import (
	"fmt"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type Args struct {
	Query     string
	Overrides parser.Result
}

type PlexName struct {
	args Args

	tmdb *tmdb.Service
	tvdb *tvdb.Service
}

func New(args Args, tmdb *tmdb.Service, tvdb *tvdb.Service) *PlexName {
	return &PlexName{
		args: args,
		tmdb: tmdb,
		tvdb: tvdb,
	}
}

func (pn *PlexName) PrintPlexName() error {
	var originalTitle string
	var err error
	parseResult := parser.Parse(pn.args.Query, pn.args.Overrides)
	if parseResult.IsMovie() {
		originalTitle, err = pn.MovieName(parseResult.Title, parseResult.Year)
	}
	if parseResult.IsTV() {
		originalTitle, err = pn.TVName(parseResult.Title, parseResult.Year)
	}
	if err != nil {
		fmt.Println("Fail: %v", err)
		return nil
	}
	fmt.Println(originalTitle)
	return nil
}
