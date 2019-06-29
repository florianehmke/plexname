package search

import (
	"fmt"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type Result struct {
	Title string
	Year  int
}

type Searcher interface {
	SearchMovie(query string, year int) ([]Result, error)
	SearchTV(query string, year int) ([]Result, error)
}

type searcher struct {
	tmdbClient tmdb.Client
	tvdbClient tvdb.Client
}

func NewSearcher(tmdbClient tmdb.Client, tvdbClient tvdb.Client) Searcher {
	return &searcher{
		tmdbClient: tmdbClient,
		tvdbClient: tvdbClient,
	}
}

func (s *searcher) SearchMovie(query string, year int) ([]Result, error) {
	response, err := s.tmdbClient.Search(query, year, 0)
	if err != nil {
		return nil, fmt.Errorf("movie search failed: %v", err)
	}
	var result []Result
	for _, r := range response.Results {
		result = append(result, Result{r.Title, r.Year()})
	}
	return result, nil
}

func (s *searcher) SearchTV(query string, year int) ([]Result, error) {
	response, err := s.tvdbClient.Search(query)
	if err != nil {
		return nil, fmt.Errorf("tv search failed: %v", err)
	}
	var result []Result
	for _, r := range response.Results {
		if year == 0 || year == r.Year() {
			result = append(result, Result{r.Title, r.Year()})
		}
	}
	return result, nil
}
