// Package tmdb provides go bindings for the TMDB API at https://www.themoviedb.org/.
package tmdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://api.themoviedb.org/3"

// Service is the TMDB service struct.
type Service struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string

	throttle chan time.Time
}

type Client interface {
	Search(query string, year int, page int) (*SearchResponse, error)
}

// NewService creates a new TMDB service.
func NewService(baseURL string, apiKey string) *Service {
	service := &Service{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
	service.startRateLimiter()
	return service
}

func (s *Service) ensureRateLimit() {
	<-s.throttle
}

func (s *Service) startRateLimiter() {
	rate := time.Second / 4
	burstLimit := 35
	ticker := time.NewTicker(rate)
	s.throttle = make(chan time.Time, burstLimit)
	go func() {
		for t := range ticker.C {
			select {
			case s.throttle <- t:
			default:
			}
		}
	}()
}

type apiError struct {
	StatusMessage string `json:"status_message"`
	StatusCode    int    `json:"status_code"`
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
		return errors.New(fmt.Sprintf("%s (code %d)", apiErr.StatusMessage, apiErr.StatusCode))
	}
	return nil
}
