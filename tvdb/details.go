package tvdb

import (
	"errors"
	"fmt"
	"strings"
)

const detailsEndpoint = "series/%s"

type DetailsResponse struct {
	DetailsResponseData `json:"data"`
}

type DetailsResponseData struct {
	Added           string   `json:"added"`
	AirsDayOfWeek   string   `json:"airsDayOfWeek"`
	AirsTime        string   `json:"airsTime"`
	Aliases         []string `json:"aliases"`
	PosterLink      string   `json:"banner"`
	FirstAired      string   `json:"firstAired"`
	GenreArray      []string `json:"genre"`
	TvdbID          int      `json:"id"`
	ImdbID          string   `json:"imdbId"`
	LastUpdated     int      `json:"lastUpdated"`
	Network         string   `json:"network"`
	NetworkID       string   `json:"networkId"`
	Plot            string   `json:"overview"`
	Rating          string   `json:"rating"`
	Runtime         string   `json:"runtime"`
	SeriesID        string   `json:"seriesId"`
	Title           string   `json:"seriesName"`
	SiteRating      float32  `json:"siteRating"`
	SiteRatingCount int      `json:"siteRatingCount"`
	Status          string   `json:"status"`
	Zap2ItID        string   `json:"zap2itId"`
}

func (dr *DetailsResponse) Genres() string {
	return strings.Join(dr.GenreArray, ", ")
}

// Details for a series fetched from TVDB.
func (s *Service) Details(seriesID string) (*DetailsResponse, error) {
	err := s.refreshTokenIfNecessary()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh jwt token: %v", err)
	}
	res := new(DetailsResponse)
	apiErr := new(apiError)
	url := fmt.Sprintf(detailsEndpoint, seriesID)
	if _, err := s.base.New().Get(url).Receive(res, apiErr); err != nil {
		return nil, fmt.Errorf("tvdb details request failed: %v", err)
	}
	if apiErr.isPresent() {
		return nil, errors.New(apiErr.Error)
	}
	return res, nil
}
