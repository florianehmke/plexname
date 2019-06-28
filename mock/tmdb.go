package mock

import "github.com/florianehmke/plexname/tmdb"

type tmdbClient struct {
	response tmdb.SearchResponse
	err      error
}

func NewMockTMDB(response tmdb.SearchResponse, err error) tmdb.Client {
	return &tmdbClient{response, err}
}

func (c *tmdbClient) Search(query string, year int, page int) (*tmdb.SearchResponse, error) {
	return &c.response, c.err
}
