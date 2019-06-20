package plexname

import (
	"errors"
	"fmt"
)

type Movie struct {
	Title         string
	Year          int
	OriginalTitle string
}

func (m *Movie) PlexName() string {
	return fmt.Sprintf("%s (%d)", m.OriginalTitle, m.Year)
}

func (pn *PlexName) MovieName(query string, year int) (string, error) {
	ot, err := pn.originalTitleFor(query, year)
	if err != nil {
		return "", fmt.Errorf("could not determine original title: %v", err)
	}
	m := Movie{Title: query, Year: year, OriginalTitle: ot}
	return m.PlexName(), nil
}

func (pn *PlexName) originalTitleFor(title string, year int) (string, error) {
	response, err := pn.tmdb.Search(title, year, 0)
	if err != nil {
		return "", fmt.Errorf("tmdb search failed: %v", err)
	}

	if len(response.Results) == 0 {
		return "", errors.New("tmdb search returned nothing")
	}
	return response.Results[0].Title, nil
}
