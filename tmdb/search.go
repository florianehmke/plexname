package tmdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SearchResponse from TMDB.
type SearchResponse struct {
	Page    int `json:"page"`
	Results []struct {
		PosterPath       string  `json:"poster_path"`
		Adult            bool    `json:"adult"`
		Plot             string  `json:"overview"`
		ReleaseDate      string  `json:"release_date"`
		GenreIds         []int   `json:"genre_ids"`
		ID               int     `json:"id"`
		OriginalTitle    string  `json:"original_title"`
		OriginalLanguage string  `json:"original_language"`
		Title            string  `json:"title"`
		BackdropPath     string  `json:"backdrop_path"`
		Popularity       float64 `json:"popularity"`
		VoteCount        int     `json:"vote_count"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
	} `json:"results"`
	TotalResults int `json:"total_results"`
	TotalPages   int `json:"total_pages"`
}

// Search for movies on TMDB.
func (s *Service) Search(query string, year int, page int) (*SearchResponse, error) {
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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read tmdb details response", err)
	}

	// Check for error.
	err = hasError(resp, b)
	if err != nil {
		return nil, err
	}

	// Unmarshal success response.
	var result SearchResponse
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal details response", err)
	}
	return &result, nil
}
