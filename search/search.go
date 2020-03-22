package search

import (
	"fmt"
	"strings"

	"github.com/florianehmke/plexname/prompt"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

type Result struct {
	Title string
	Year  int
}

type Query struct {
	Title string
	Year  int
}

type Searcher interface {
	SearchMovie(query Query) (Result, error)
	SearchTV(query Query) (Result, error)
}

type searcher struct {
	tmdbClient tmdb.Client
	tvdbClient tvdb.Client
	prompter   prompt.Prompter

	cache map[Query]Result
}

func NewSearcher(tmdbClient tmdb.Client, tvdbClient tvdb.Client, prompter prompt.Prompter) Searcher {
	return &searcher{
		tmdbClient: tmdbClient,
		tvdbClient: tvdbClient,
		prompter:   prompter,
		cache:      map[Query]Result{},
	}
}

func (s *searcher) SearchMovie(query Query) (Result, error) {
	if v, ok := s.cache[query]; ok {
		return v, nil
	}
	response, err := s.tmdbClient.Search(query.Title, query.Year, 0)
	if err != nil {
		return Result{}, fmt.Errorf("movie search failed: %v", err)
	}
	var result []Result
	for _, r := range response.Results {
		result = append(result, Result{r.Title, r.Year()})
	}
	if len(result) == 0 {
		fmt.Printf("no search result for title '%s'\n", query.Title)
		query, err := s.prompter.AskString(fmt.Sprintf("Search again:"))
		if err != nil {
			return Result{}, fmt.Errorf("prompt error: %v", err)
		}
		return s.SearchMovie(Query{Title: query})
	}
	return s.toSingleResult(query, result)
}

func (s *searcher) SearchTV(query Query) (Result, error) {
	if v, ok := s.cache[query]; ok {
		return v, nil
	}
	response, err := s.tvdbClient.Search(query.Title)
	if err != nil {
		return Result{}, fmt.Errorf("tv search failed: %v", err)
	}
	var result []Result
	for _, r := range response.Results {
		result = append(result, Result{r.Title, r.Year()})
	}
	if len(result) == 0 {
		fmt.Printf("no search result for title '%s'\n", query.Title)
		query, err := s.prompter.AskString(fmt.Sprintf("Search again:"))
		if err != nil {
			return Result{}, fmt.Errorf("prompt error: %v", err)
		}
		return s.SearchTV(Query{Title: query})
	}
	return s.toSingleResult(query, result)
}

func (s *searcher) toSingleResult(query Query, results []Result) (Result, error) {
	var result Result
	if len(results) > 1 {
		choices := []string{fmt.Sprintf("Multiple results found online for %s, pick one of:", query.Title)}
		for i, r := range results {
			choices = append(choices, fmt.Sprintf("[%d] %s (%d)", i+1, r.Title, r.Year))
		}
		i, err := s.prompter.AskNumber(strings.Join(choices, "\n"))
		if err != nil {
			return result, fmt.Errorf("prompt error: %v", err)
		}
		result = results[i-1]
	} else {
		result = results[0]
	}
	s.cache[query] = result
	return result, nil
}
