// Package tmdb provides go bindings for the TMDB API at https://www.themoviedb.org/.
package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"errors"
)

const (
	BaseURL        = "https://api.themoviedb.org/3"
	searchEndpoint = "/search/%s?api_key=%s&language=en-US"
	movieEndpoint  = "/movie/%d?api_key=%s&language=en-US"
)

// Service is the TMDB service struct.
type Service struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string

	throttle chan time.Time
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

func hasError(resp *http.Response, body []byte) error {
	if resp.StatusCode != http.StatusOK {
		var error apiError
		err := json.Unmarshal(body, &error)
		if err != nil {
			return fmt.Errorf("could not unmarshal error response: %v", err)
		}
		return errors.New(error.StatusMessage)
	}
	return nil
}
