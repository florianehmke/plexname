package tvdb

import (
	"errors"
	"fmt"
)

const searchEndpoint = "search/series?name=%s"

type SearchResponse struct {
	Results []struct {
		Aliases    []string `json:"aliases"`
		PosterLink string   `json:"banner"`
		FirstAired string   `json:"firstAired"`
		TvdbID     int      `json:"id"`
		Network    string   `json:"network"`
		Plot       string   `json:"overview"`
		Title      string   `json:"seriesName"`
		Status     string   `json:"status"`
	} `json:"data"`
}

// Search for series on TVDB.
func (s *Service) Search(query string) (*SearchResponse, error) {
	err := s.refreshTokenIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh jwt token: %v", err)
	}
	res := new(SearchResponse)
	apiErr := new(apiError)
	url := fmt.Sprintf(searchEndpoint, query)
	if _, err := s.base.New().Get(url).Receive(res, apiErr); err != nil {
		return nil, fmt.Errorf("tvdb search request failed: %v", err)
	}
	if apiErr.isPresent() {
		return nil, errors.New(apiErr.Error)
	}
	return res, nil
}
