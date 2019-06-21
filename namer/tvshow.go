package namer

import (
	"errors"
	"fmt"
)

func (pn *Namer) originalTvShowTitleFor(title string, year int) (string, error) {
	response, err := pn.tvdb.Search(title)
	if err != nil {
		return "", fmt.Errorf("tmdb search failed: %v", err)
	}
	if len(response.Results) == 0 {
		return "", errors.New("tmdb search returned nothing")
	}
	return response.Results[0].Title, nil
}
