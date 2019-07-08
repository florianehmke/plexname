package renamer_test

import (
	"path/filepath"
	"testing"

	"github.com/florianehmke/plexname/mock"
	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/renamer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type testFixture struct {
	renamer.Parameters

	tvdbResponse []tvdb.SearchResult
	tmdbResponse []tmdb.SearchResult

	promptResponse int

	expectedOldFilePath string
	expectedNewFilePath string
	expectedNewPath     string
}

var tests = []testFixture{
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/movie-parse-from-file",
			TargetPath: "/dev/null/",
			OnlyFile:   true,
		},
		expectedOldFilePath: "../tests/fixtures/movie-parse-from-file/movie title/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title (1999)/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title (1999)/",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		Parameters: renamer.Parameters{

			SourcePath: "../tests/fixtures/movie-parse-from-folder",
			TargetPath: "/dev/null/",
			OnlyDir:    true,
		},
		expectedOldFilePath: "../tests/fixtures/movie-parse-from-folder/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group/movie.file.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title (1999)/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title (1999)/",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/movie-with-tmdb-prompt",
			TargetPath: "/dev/null/",
		},
		expectedOldFilePath: "../tests/fixtures/movie-with-tmdb-prompt/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group/movie.file.mkv",
		expectedNewFilePath: "/dev/null/Real Movie Title 2 (1999)/Real Movie Title 2 (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real Movie Title 2 (1999)/",
		promptResponse:      2,
		tmdbResponse: []tmdb.SearchResult{
			{Title: "Real Movie Title 1"},
			{Title: "Real Movie Title 2"},
		},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-parse-from-file",
			TargetPath: "/dev/null/",
		},
		expectedOldFilePath: "../tests/fixtures/tv-parse-from-file/tv show title/TV-Show.S02E13.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Real TV Show Title/Season 02/Real TV Show Title - S02E13 - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Real TV Show Title/Season 02/",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Real TV Show Title"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-with-tvdb-prompt",
			TargetPath: "/dev/null/",
		},
		expectedOldFilePath: "../tests/fixtures/tv-with-tvdb-prompt/tv show title/TV-Show.S02E13.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "/dev/null/Another Real TV Show Title (1981)/Season 02/Another Real TV Show Title (1981) - S02E13 - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "/dev/null/Another Real TV Show Title (1981)/Season 02/",
		promptResponse:      2,
		tvdbResponse: []tvdb.SearchResult{
			{Title: "Real TV Show Title"},
			{Title: "Another Real TV Show Title", FirstAired: "1981-01-01"},
		},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/movie-file-only/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		},
		expectedOldFilePath: "../tests/fixtures/movie-file-only/Movie.Title.1999.German.1080p.DL.DTS.BluRay.AVC.Remux-group.mkv",
		expectedNewFilePath: "../tests/fixtures/movie-file-only/Real Movie Title (1999) - German.1080p.DL.Blu-ray.Remux.mkv",
		expectedNewPath:     "../tests/fixtures/movie-file-only/",
		tmdbResponse:        []tmdb.SearchResult{{Title: "Real Movie Title"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-standalone-episode-numbers",
		},
		expectedOldFilePath: "../tests/fixtures/tv-standalone-episode-numbers/Awesome Show S01/1 - Title.mkv",
		expectedNewFilePath: "../tests/fixtures/tv-standalone-episode-numbers/Awesome Show/Season 01/Awesome Show - S01E01.mkv",
		expectedNewPath:     "../tests/fixtures/tv-standalone-episode-numbers/Awesome Show/Season 01/",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Awesome Show"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-dual-ep",
		},
		expectedOldFilePath: "../tests/fixtures/tv-dual-ep/tv show title/Title S01E03E04.mkv",
		expectedNewFilePath: "../tests/fixtures/tv-dual-ep/Awesome Show/Season 01/Awesome Show - S01E03E04.mkv",
		expectedNewPath:     "../tests/fixtures/tv-dual-ep/Awesome Show/Season 01/",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Awesome Show"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-long-season-folder-name",
		},
		expectedOldFilePath: "../tests/fixtures/tv-long-season-folder-name/A.Very.Long.Show.Name.2000.German.S01.DL.1080p.BluRay/Title.E06.mkv",
		expectedNewFilePath: "../tests/fixtures/tv-long-season-folder-name/Awesome Show (2000)/Season 01/Awesome Show (2000) - S01E06 - German.1080p.DL.Blu-ray.mkv",
		expectedNewPath:     "../tests/fixtures/tv-long-season-folder-name/Awesome Show (2000)/Season 01/",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Awesome Show", FirstAired: "2000-01-01"}},
	},
	{
		Parameters: renamer.Parameters{
			SourcePath: "../tests/fixtures/tv-special",
		},
		expectedOldFilePath: "../tests/fixtures/tv-special/tv show title/specials/Special E01 German.mkv",
		expectedNewFilePath: "../tests/fixtures/tv-special/Awesome Show/Season 00/Awesome Show - S00E01 - German.mkv",
		expectedNewPath:     "../tests/fixtures/tv-special/Awesome Show/Season 00/",
		tvdbResponse:        []tvdb.SearchResult{{Title: "Awesome Show"}},
	},
}

func TestFixtures(t *testing.T) {
	for _, tc := range tests {
		mockedTVDB := mockTVDBResponse(tc.tvdbResponse)
		mockedTMDB := mockTMDBResponse(tc.tmdbResponse)
		mockedFS := mock.NewMockFS(func(oldPath string, newPath string) error {
			if oldPath != tc.expectedOldFilePath {
				t.Errorf("\nExpected: %s\nReceived: %s", tc.expectedOldFilePath, oldPath)
			}
			if newPath != tc.expectedNewFilePath {
				t.Errorf("\nExpected: %s\nReceived: %s", tc.expectedNewFilePath, newPath)
			}
			return nil
		}, func(path string) error {
			if path != tc.expectedNewPath {
				t.Errorf("\nExpected: %s\nReceived: %s", tc.expectedNewPath, path)
			}
			return nil
		})
		mockedPrompter := mock.NewMockPrompter(func(question string) (i int, e error) {
			return tc.promptResponse, nil
		}, nil, nil)

		sourcePath := filepath.FromSlash(tc.SourcePath)
		targetPath := filepath.FromSlash(tc.TargetPath)
		n := renamer.New(
			renamer.NewParameters(sourcePath, targetPath, parser.Result{}, []string{}, false, false, false),
			search.NewSearcher(mockedTMDB, mockedTVDB, mockedPrompter),
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
