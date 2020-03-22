// Package tvdb provides go bindings for the TVDB API at https://api.thetvdb.com/swagger.
package tvdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://api.thetvdb.com/"

// client is the TVDB client struct.
type client struct {
	client *http.Client

	apiKey  string
	baseURL string

	token         tokenResponse
	tokenFromDate time.Time
}

type Client interface {
	Search(query string) (*SearchResponse, error)
}

// New creates a new TMDB client.
func NewClient(baseURL string, apiKey string) Client {
	tvdbService := &client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
		token:   tokenResponse{},
	}
	return tvdbService
}

func (s *client) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))
	req.Header.Add("Content-Type", "application/json")
}

type apiError struct {
	Error string `json:"Error"`
}

func unmarshalResponse(resp *http.Response, success interface{}) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if success != nil && resp.StatusCode != 204 {
			return json.NewDecoder(resp.Body).Decode(success)
		}
	} else {
		var apiErr apiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("unmarshal of response failed: %v", err)
		}
		return errors.New(apiErr.Error)
	}
	return nil
}
