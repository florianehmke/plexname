package plexname

import (
	"errors"
	"fmt"
)

type TV struct {
	Title         string
	Episode int
	Season int
	OriginalTitle string
}

func (tv *TV) PlexName() string {
	return fmt.Sprintf("%v", tv)
}

func (pn *PlexName) TVName(query string, year int) (string, error) {
	ot, err := pn.originalTvShowTitleFor(query, year)
	if err != nil {
		return "", fmt.Errorf("could not determine original title: %v", err)
	}
	m := Movie{Title: query, Year: year, OriginalTitle: ot}
	return m.PlexName(), nil
}

func (pn *PlexName) originalTvShowTitleFor(title string, year int) (string, error) {
	response, err := pn.tvdb.Search(title)
	if err != nil {
		return "", fmt.Errorf("tmdb search failed: %v", err)
	}

	if len(response.Results) == 0 {
		return "", errors.New("tmdb search returned nothing")
	}
	return response.Results[0].Title, nil
}
