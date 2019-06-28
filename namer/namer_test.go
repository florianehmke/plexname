package namer_test

import (
	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
	"testing"
)

func TestParseFromFile(t *testing.T) {
	mockedTVDB := newMockTVDB("hi", nil)
	mockedTMDB := newMockTMDB("hi", nil)
	fs := mock.NewMockFS(nil)

	n := namer.New(namer.Args{Path: "../tests/fixtures/parse-from-file"}, mockedTMDB, mockedTVDB, fs)
	if err := n.Run(); err != nil {
		t.Error(err)
	}
}

func newMockTVDB(result string, err error) tvdb.Client {
	return mock.NewMockTVDB(tvdb.SearchResponse{
		Results: []tvdb.SearchResult{
			{Title: result},
		},
	}, err)
}

func newMockTMDB(result string, err error) tmdb.Client {
	return mock.NewMockTMDB(tmdb.SearchResponse{
		Results: []tmdb.SearchResult{
			{Title: result},
		},
	}, err)
}
