package mock

import "github.com/florianehmke/plexname/tvdb"

type tvdbClient struct {
	response tvdb.SearchResponse
	err      error
}

func NewMockTVDB(response tvdb.SearchResponse, err error) tvdb.Client {
	return &tvdbClient{response, err}
}

func (c *tvdbClient) Search(query string) (*tvdb.SearchResponse, error) {
	return &c.response, c.err
}
