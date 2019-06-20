// Package tvdb provides go bindings for the TVDB API at https://api.thetvdb.com/swagger.
package tvdb

import (
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/florianehmke/plexname/log"
)

const BaseURL = "https://api.thetvdb.com/"

// Service is the TVDB service struct.
type Service struct {
	base *sling.Sling

	apiKey string

	token         token
	tokenFromDate time.Time
}

type token struct {
	JWTToken string `json:"token"`
}

type authRequestBody struct {
	Apikey string `json:"apikey"`
	// not used..
	Userkey  string `json:"userkey"`
	Username string `json:"username"`
}

type apiError struct {
	Error string `json:"Error"`
}

func (e apiError) isPresent() bool {
	// Check if an empty error equals e
	return (apiError{}) != e
}

// New creates a new TMDB service.
func NewService(baseURL string, apiKey string) *Service {
	tvdbService := &Service{
		apiKey: apiKey,
		base:   sling.New().Base(baseURL),
	}
	return tvdbService
}

// authenticate at TVDB.
func (s *Service) authenticate() error {
	tok := new(token)
	apiErr := new(apiError)
	body := authRequestBody{Apikey: s.apiKey}
	if _, err := s.base.New().Post("/login").BodyJSON(body).Receive(tok, apiErr); err != nil {
		return fmt.Errorf("tvdb authentication request failed: %v", err)
	}
	if apiErr.isPresent() {
		return fmt.Errorf("tvdb authentication failed: %s", apiErr.Error)
	}
	s.token = *tok
	s.tokenFromDate = time.Now()
	s.base.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))
	log.Infof("Received tvdb token: %s...", s.token.JWTToken[:10])
	return nil
}

// refreshToken at TVDB.
func (s *Service) refreshToken() error {
	tok := new(token)
	apiErr := new(apiError)
	if _, err := s.base.New().Get("/login").Receive(tok, apiErr); err != nil {
		return fmt.Errorf("tvdb authentication request failed: %v", err)
	}
	if apiErr.isPresent() {
		return fmt.Errorf("tvdb authentication failed: %s", apiErr.Error)
	}
	s.base.Set("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))
	return nil
}

// check if a jwt token refresh is necessary and do it if so.
func (s *Service) refreshTokenIfNecessary() error {
	dur := time.Since(s.tokenFromDate)
	if 18 < dur.Hours() && dur.Hours() < 24 {
		log.Info("Refreshing tvdb token..")
		return s.refreshToken()
	}
	if dur.Hours() > 24 || s.token.JWTToken == "" {
		log.Info("Authenticating with tvdb.")
		return s.authenticate()
	}
	return nil
}
