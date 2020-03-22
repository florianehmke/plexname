package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/florianehmke/plexname/config"
)

const authEndpoint = "login"

type tokenResponse struct {
	JWTToken string `json:"token"`
}

type authRequestBody struct {
	Apikey string `json:"apikey"`
}

func (s *client) requestInitialToken() error {
	body, err := json.Marshal(authRequestBody{Apikey: config.GetToken("tvdb")})
	if err != nil {
		return fmt.Errorf("marshal of request body failed: %v", err)
	}
	resp, err := http.Post(s.baseURL+authEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http post failed: %v", err)
	}
	defer resp.Body.Close()
	if err := unmarshalResponse(resp, &s.token); err != nil {
		return fmt.Errorf("unmarshal of response failed: %v", err)
	}
	return nil
}

func (s *client) refreshToken() error {
	req, err := http.NewRequest("GET", s.baseURL+authEndpoint, nil)
	if err != nil {
		return fmt.Errorf("creation of get request failed: %v", err)
	}
	s.addHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("http get failed: %v", err)
	}
	defer resp.Body.Close()

	if err := unmarshalResponse(resp, &s.token); err != nil {
		return fmt.Errorf("unmarshal of response failed: %v", err)
	}
	return nil
}

func (s *client) refreshTokenIfNecessary() error {
	dur := time.Since(s.tokenFromDate)
	if 18 < dur.Hours() && dur.Hours() < 24 {
		return s.refreshToken()
	}
	if dur.Hours() > 24 || s.token.JWTToken == "" {
		return s.requestInitialToken()
	}
	return nil
}
