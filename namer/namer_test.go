package namer_test

import (
	"path/filepath"
	"testing"

	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/namer"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type testFixture struct {
	sourcePath string
	targetPath string

	tvdbResponse []tvdb.SearchResult
	tmdbResponse []tmdb.SearchResult

	promptResponse int

	expectedOldFilePath string
	expectedNewFilePath string
	expectedNewPath     string
}

var tests = []testFixture{
	{
		sourcePath:          "../tests/fixtures/movie-parse-from-file",
		targetPath:          "/dev/null",
		expectedOldFilePath: "../tests/fixtures/movie-parse-from-file/movie title/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title (1999)/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title (1999)",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		sourcePath:          "../tests/fixtures/movie-parse-from-folder",
		targetPath:          "/dev/null",
		expectedOldFilePath: "../tests/fixtures/movie-parse-from-folder/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group/movie.file.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title (1999)/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title (1999)",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		sourcePath:          "../tests/fixtures/movie-with-tmdb-prompt",
		targetPath:          "/dev/null",
		expectedOldFilePath: "../tests/fixtures/movie-with-tmdb-prompt/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group/movie.file.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title 2 (1999)/Real Movie Title 2 (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title 2 (1999)",
		promptResponse:      2,
		tmdbResponse: []tmdb.SearchResult{
			{Title: "Real Movie Title 1"},
			{Title: "Real Movie Title 2"},
		},
	},
	{
		sourcePath:          "../tests/fixtures/tv-parse-from-file",
		targetPath:          "/dev/null",
		expectedOldFilePath: "../tests/fixtures/tv-parse-from-file/tv show title/TV-Show.S02E13.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Real TV Show Title/Season 02/Real TV Show Title - s02e13 - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real TV Show Title/Season 02",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Real TV Show Title"}},
	},
	{
		sourcePath:          "../tests/fixtures/tv-with-tvdb-prompt",
		targetPath:          "/dev/null",
		expectedOldFilePath: "../tests/fixtures/tv-with-tvdb-prompt/tv show title/TV-Show.S02E13.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Another Real TV Show Title (1981)/Season 02/Another Real TV Show Title (1981) - s02e13 - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Another Real TV Show Title (1981)/Season 02",
		promptResponse:      2,
		tvdbResponse: []tvdb.SearchResult{
			{Title: "Real TV Show Title"},
			{Title: "Another Real TV Show Title", FirstAired: "1981-01-01"},
		},
	},
	{
		sourcePath:          "../tests/fixtures/movie-file-only/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedOldFilePath: "../tests/fixtures/movie-file-only/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "../tests/fixtures/movie-file-only/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
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

		sourcePath := filepath.FromSlash(tc.sourcePath)
		targetPath := filepath.FromSlash(tc.targetPath)
		n := namer.New(
			namer.NewArgs(sourcePath, targetPath, parser.Result{}),
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
