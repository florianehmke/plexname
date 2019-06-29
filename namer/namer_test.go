package namer_test

import (
	"log"
	"testing"

	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func TestParseFromFile(t *testing.T) {
	mockedTVDB := mockTVDBResponse("some tvshow", "")
	mockedTMDB := mockTMDBResponse("some movie", "")
	mockedFS := mock.NewMockFS(func(oldPath string, newPath string) error {
		log.Println(oldPath, newPath)
		return nil
	}, func(path string) error {
		log.Println(path)
		return nil
	})

	n := namer.New(
		namer.Args{Path: "../tests/fixtures/parse-from-file"},
		search.NewSearcher(mockedTMDB, mockedTVDB), mockedFS)
	if err := n.Run(); err != nil {
		t.Error(err)
	}
}

func mockTVDBResponse(result string, firstAired string) tvdb.Client {
	return mock.NewMockTVDB(tvdb.SearchResponse{
		Results: []tvdb.SearchResult{{
			Title:      result,
			FirstAired: firstAired,
		}},
	}, nil)
}

func mockTMDBResponse(result string, releaseDate string) tmdb.Client {
	return mock.NewMockTMDB(tmdb.SearchResponse{
		Results: []tmdb.SearchResult{{
			Title:       result,
			ReleaseDate: releaseDate,
		}},
	}, nil)
}
