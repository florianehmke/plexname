package namer_test

import (
	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
	"testing"
)

func TestParseFromFile(t *testing.T) {
	tvdbResponse := &tvdb.SearchResponse{}
	var tvdbError error
	mockedTVDB := mock.NewMockTVDB(tvdbResponse, &tvdbError)

	tmdbResponse := &tmdb.SearchResponse{Results: []tmdb.SearchResult{{Title: "hi"}}}
	var tmdbError error
	mockedTMDB := mock.NewMockTMDB(tmdbResponse, &tmdbError)

	var fsError error
	fs := mock.NewMockFS(&fsError)

	n := namer.New(namer.Args{Path: "../tests/fixtures/parse-from-file"}, mockedTMDB, mockedTVDB, fs)
	if err := n.Run(); err != nil {
		t.Error(err)
	}
}
