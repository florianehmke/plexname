package tmdb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const searchEndpoint = "/search/%s?api_key=%s&language=en-US"

// SearchResponse from TMDB.
type SearchResponse struct {
	Page         int            `json:"page"`
	Results      []SearchResult `json:"results"`
	TotalResults int            `json:"total_results"`
	TotalPages   int            `json:"total_pages"`
}

type SearchResult struct {
	ReleaseDate string `json:"release_date"` // e.g. 2014-03-20
	Title       string `json:"title"`
}

func (sr *SearchResult) Year() int {
	year := 0
	if sr.ReleaseDate != "" && len(sr.ReleaseDate) >= 4 {
		yearString := sr.ReleaseDate[:4]
		if y, err := strconv.Atoi(yearString); err == nil {
			year = y
		}
	}
	return year
}

// Search for movies on TMDB.
func (s *client) Search(query string, year int, page int) (*SearchResponse, error) {
	reqURL := fmt.Sprintf(s.baseURL+searchEndpoint, "movie", s.apiKey)

	// Build the query string.
	v := url.Values{}
	v.Set("query", query)
	if year > 0 {
		v.Add("year", strconv.Itoa(year))
	}
	if page > 0 {
		v.Add("page", strconv.Itoa(page))
	}
	qs := v.Encode()
	if qs != "" {
		reqURL = reqURL + "&" + qs
	}

	// Do the request.
	s.ensureRateLimit()
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("could not get tmdb details: %v", err)
	}
	defer resp.Body.Close()

	var result SearchResponse
	if err := unmarshalResponse(resp, &result); err != nil {
		return nil, fmt.Errorf("unmarshal of response failed: %v", err)
	}
	return &result, nil
}
