package namer_test

import (
	"path/filepath"
	"testing"

	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type testFixture struct {
	fixturePath string

	tvdbResponse []tvdb.SearchResult
	tmdbResponse []tmdb.SearchResult

	promptResponse int

	expectedOldFilePath string
	expectedNewFilePath string
	expectedNewPath     string
}

var tests = []testFixture{
	{
		fixturePath:         "../tests/fixtures/parse-from-file",
		expectedOldFilePath: "../tests/fixtures/parse-from-file/movie title/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "../tests/fixtures/parse-from-file/Real Movie Title (1999)/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewPath:     "../tests/fixtures/parse-from-file/Real Movie Title (1999)",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		fixturePath:         "../tests/fixtures/parse-from-folder",
		expectedOldFilePath: "../tests/fixtures/parse-from-folder/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group/movie.file.mkv",
		expectedNewFilePath: "../tests/fixtures/parse-from-folder/Real Movie Title (1999)/movie.file.mkv",
		expectedNewPath:     "../tests/fixtures/parse-from-folder/Real Movie Title (1999)",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
}

func TestFixtures(t *testing.T) {
	for _, tc := range tests {
		mockedTVDB := mockTVDBResponse(tc.tvdbResponse)
		mockedTMDB := mockTMDBResponse(tc.tmdbResponse)
		mockedFS := mock.NewMockFS(func(oldPath string, newPath string) error {
			if oldPath != tc.expectedOldFilePath {
				t.Errorf("expected %s, got %s", tc.expectedOldFilePath, oldPath)
			}
			if newPath != tc.expectedNewFilePath {
				t.Errorf("expected %s, got %s", tc.expectedNewFilePath, newPath)
			}
			return nil
		}, func(path string) error {
			if path != tc.expectedNewPath {
				t.Errorf("expected %s, got %s", tc.expectedNewPath, path)
			}
			return nil
		})
		mockedPrompter := mock.NewMockPrompter(func(question string) (i int, e error) {
			return tc.promptResponse, nil
		}, nil)

		n := namer.New(
			namer.Args{Path: filepath.FromSlash(tc.fixturePath)},
			search.NewSearcher(mockedTMDB, mockedTVDB),
			mockedPrompter,
			mockedFS)
		if err := n.Run(); err != nil {
			t.Error(err)
		}
	}
}

func mockTVDBResponse(results []tvdb.SearchResult) tvdb.Client {
	return mock.NewMockTVDB(tvdb.SearchResponse{Results: results}, nil)
}

func mockTMDBResponse(results []tmdb.SearchResult) tmdb.Client {
	return mock.NewMockTMDB(tmdb.SearchResponse{Results: results}, nil)
}
