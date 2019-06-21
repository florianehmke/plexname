// Package tvdb provides go bindings for the TVDB API at https://api.thetvdb.com/swagger.
package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/log"
)

const BaseURL = "https://api.thetvdb.com/"

// Service is the TVDB service struct.
type Service struct {
	client *http.Client

	apiKey string

	token         token
	tokenFromDate time.Time
}

type token struct {
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
	}
	return tvdbService
}

// authenticate at TVDB.
func (s *Service) authenticate() error {
	body, err := json.Marshal(authRequestBody{Apikey: config.GetToken("tvdb")})
	if err != nil {
		return fmt.Errorf("could not marshal auth request body: %v", err)
	}
	resp, err := http.Post(BaseURL+"login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("could not post auth request: %v", err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if err := json.NewDecoder(resp.Body).Decode(&s.token); err != nil {
			return fmt.Errorf("could not unmarshal auth response: %v", err)
		}
	} else {
		var apiErr apiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("could not unmarshal auth error response: %v", err)
		}
		return fmt.Errorf("tvdb authentication failed: %s", apiErr.Error)
	}
	log.Infof("Received tvdb token: %s...", s.token.JWTToken[:10])
	return nil
}

// refreshToken at TVDB.
func (s *Service) refreshToken() error {
	req, err := http.NewRequest("GET", BaseURL+"login", nil)
	if err != nil {
		return fmt.Errorf("could not create refresh token request: %v", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token.JWTToken))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("coult not do refresh token request: %v", err)
	}
	defer resp.Body.Close()

	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if err := json.NewDecoder(resp.Body).Decode(&s.token); err != nil {
			return fmt.Errorf("could not unmarshal auth response: %v", err)
		}
	} else {
		var apiErr apiError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("could not unmarshal auth error response: %v", err)
		}
		return fmt.Errorf("tvdb authentication failed: %s", apiErr.Error)
	}
	log.Infof("Refreshed tvdb token: %s...", s.token.JWTToken[:10])
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
