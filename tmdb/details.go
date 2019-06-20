package tmdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DetailsResponse from TMDB.
type DetailsResponse struct {
	Adult               bool        `json:"adult"`
	BackdropPath        string      `json:"backdrop_path"`
	BelongsToCollection interface{} `json:"belongs_to_collection"`
	Budget              int         `json:"budget"`
	GenreArray          []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	TMDBID              int     `json:"id"`
	IMDBID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Plot                string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ProductionCompanies []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		Iso6391 string `json:"iso_639_1"`
		Name    string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// Genres returns the genres as a joined string.
func (dp *DetailsResponse) Genres() string {
	genres := make([]string, len(dp.GenreArray))
	for k, genre := range dp.GenreArray {
		genres[k] = genre.Name
	}
	return strings.Join(genres, ", ")
}

// Details for a movie with the given movie ID.
func (s *Service) Details(id int) (*DetailsResponse, error) {
	url := fmt.Sprintf(s.baseURL+movieEndpoint, id, s.apiKey)

	// Do the request.
	s.ensureRateLimit()
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get tmdb details: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read tmdb details response: %v", err)
	}

	// Check for error.
	err = hasError(resp, b)
	if err != nil {
		return nil, err
	}

	// Unmarshal success response.
	var result DetailsResponse
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal details response: %v", err)
	}
	return &result, nil
}
