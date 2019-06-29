package namer_test

import (
	"testing"

	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func TestParseFromFile(t *testing.T) {
	mockedTVDB := mockTVDBResponse("", "")
	mockedTMDB := mockTMDBResponse("Real Movie Title", "")
	mockedFS := mock.NewMockFS(func(oldPath string, newPath string) error {
		expectedOldPath := "../tests/fixtures/parse-from-file/movie title/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv"
		if oldPath != expectedOldPath {
			t.Errorf("expected %s, got %s", expectedOldPath, oldPath)
		}
		expectedNewPath := "../tests/fixtures/parse-from-file/Real Movie Title (1999)/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv"
		if newPath != expectedNewPath {
			t.Errorf("expected %s, got %s", expectedNewPath, newPath)
		}
		return nil
	}, func(path string) error {
		expectedPath := "../tests/fixtures/parse-from-file/Real Movie Title (1999)"
		if path != expectedPath {
			t.Errorf("expected %s, got %s", expectedPath, path)
		}
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
