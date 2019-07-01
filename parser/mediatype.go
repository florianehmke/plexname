package parser

import (
	"fmt"
	"strings"
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

func ParseMediaTypeFromPath(path string) MediaType {
	joined := strings.ToLower(clean(path))
	tokens := tokenize(strings.ToLower(path))

	for _, t := range tokens {
		if mt, ok := mediaTypes[t]; ok {
			return mt
		}
	}
	for k, mt := range mediaTypes {
		if len(k) < 5 {
			continue
		}
		if strings.Contains(joined, k) {
			return mt
		}
	}
	return MediaTypeUnknown
}

func ParseMediaType(s string) (MediaType, error) {
	if mt, ok := mediaTypes[s]; ok {
		return mt, nil
	}
	if s == "" {
		return MediaTypeUnknown, nil
	}
	return MediaTypeUnknown, fmt.Errorf("unknown media type: %s", s)
}
