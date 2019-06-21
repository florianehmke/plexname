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

// Service is the TVDB service struct.
type Service struct {
	client *http.Client

	apiKey string

	token         tokenResponse
	tokenFromDate time.Time
}

type tokenResponse struct {
	JWTToken string `json:"token"`
}

type authRequestBody struct {
	Apikey string `json:"apikey"`
}

type apiError struct {
	Error string `json:"Error"`
}

// New creates a new TMDB service.
func NewService(baseURL string, apiKey string) *Service {
	tvdbService := &Service{
		apiKey: apiKey,
		client: &http.Client{},
		token:  tokenResponse{},
	}
	return tvdbService
}

func (s *Service) addHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))
	req.Header.Add("Content-Type", "application/json")
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
