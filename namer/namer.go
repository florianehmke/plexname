package namer

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

type Namer struct {
	args Args

	tmdb *tmdb.Service
	tvdb *tvdb.Service
}

func New(args Args, tmdb *tmdb.Service, tvdb *tvdb.Service) *Namer {
	return &Namer{
		args: args,
		tmdb: tmdb,
		tvdb: tvdb,
	}
}

func (pn *Namer) PrintPlexName() error {
	var originalTitle string
	var err error
	parseResult := parser.Parse(pn.args.Query, pn.args.Overrides)
	if parseResult.IsMovie() {
		originalTitle, err = pn.originalMovieTitleFor(parseResult.Title, parseResult.Year)
	}
	if parseResult.IsTV() {
		originalTitle, err = pn.originalTvShowTitleFor(parseResult.Title, parseResult.Year)
	}
	if err != nil {
		fmt.Printf("Fail: %v\n", err)
		return nil
	}
	fmt.Println(originalTitle)
	return nil
}
