package tvdb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const searchEndpoint = "search/series?name=%s"

type SearchResponse struct {
	Results []SearchResult `json:"data"`
}

type SearchResult struct {
	FirstAired string `json:"firstAired"` // e.g. 1981-01-01
	Title      string `json:"seriesName"`
}

func (sr *SearchResult) Year() int {
	year := 0
	if sr.FirstAired != "" && len(sr.FirstAired) >= 4 {
		yearString := sr.FirstAired[:4]
		if y, err := strconv.Atoi(yearString); err == nil {
			year = y
		}
	}
	return year
}

// Search for series on TVDB.
func (s *client) Search(query string) (*SearchResponse, error) {
	err := s.refreshTokenIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("jwt token refresh failed: %v", err)
	}

	reqURL := fmt.Sprintf(s.baseURL+searchEndpoint, url.QueryEscape(query))
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
