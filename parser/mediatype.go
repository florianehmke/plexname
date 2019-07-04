package parser

import (
	"fmt"
)

type MediaType int

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeMovie
	MediaTypeTV
)

var (
	mediaTypes = map[string]MediaType{
		"tv":     MediaTypeTV,
		"series": MediaTypeTV,
		"shows":  MediaTypeTV,
		"movie":  MediaTypeMovie,
		"movies": MediaTypeMovie,
		"filme":  MediaTypeMovie,
	}
)

func ParseMediaType(s string) (MediaType, error) {
	if mt, ok := mediaTypes[s]; ok {
		return mt, nil
	}
	if s == "" {
		return MediaTypeUnknown, nil
	}
	return MediaTypeUnknown, fmt.Errorf("unknown media type: %s", s)
}
