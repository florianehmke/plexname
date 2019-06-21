package tvdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/florianehmke/plexname/log"
)

const searchEndpoint = "search/series?name=%s"

type SearchResponse struct {
	Results []struct {
		FirstAired string `json:"firstAired"`
		Title      string `json:"seriesName"`
	} `json:"data"`
}

// Search for series on TVDB.
func (s *Service) Search(query string) (*SearchResponse, error) {
	err := s.refreshTokenIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh jwt token: %v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(BaseURL+searchEndpoint, url.QueryEscape(query)), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create refresh token request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("coult not do refresh token request: %v", err)
	}
	defer resp.Body.Close()

	var result SearchResponse
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("could not unmarshal auth response: %v", err)
		}
	} else {
		var apiErr apiError
		b, _ := ioutil.ReadAll(resp.Body)
		log.Infof("%s", string(b))
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("could not unmarshal auth error response: %v", err)
		}
		return nil, fmt.Errorf("tvdb authentication failed: %s", apiErr.Error)
	}
	return &result, nil
}
