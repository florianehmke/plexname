package tvdb

import (
	"fmt"
	"net/http"
	"net/url"
)

const searchEndpoint = "search/series?name=%s"

type SearchResponse struct {
	Results []SearchResult `json:"data"`
}

type SearchResult struct {
	FirstAired string `json:"firstAired"`
	Title      string `json:"seriesName"`
}

// Search for series on TVDB.
func (s *client) Search(query string) (*SearchResponse, error) {
	err := s.refreshTokenIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("jwt token refresh failed: %v", err)
	}

	reqURL := fmt.Sprintf(BaseURL+searchEndpoint, url.QueryEscape(query))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creation of get request failed: %v", err)
	}
	s.addHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %v", err)
	}
	defer resp.Body.Close()

	var result SearchResponse
	if err := unmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("unmarshal of response failed: %v", err)
	}
	return &result, nil
}
