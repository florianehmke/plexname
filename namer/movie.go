package namer

import (
	"errors"
	"fmt"
)

func (pn *Namer) originalMovieTitleFor(title string, year int) (string, error) {
	response, err := pn.tmdb.Search(title, year, 0)
	if err != nil {
		return "", fmt.Errorf("tmdb search failed: %v", err)
	}
	if len(response.Results) == 0 {
		return "", errors.New("tmdb search returned nothing")
	}
	return response.Results[0].Title, nil
}
