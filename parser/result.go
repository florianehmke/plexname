package parser

import "strings"

type Result struct {
	Title string

	MediaType MediaType

	Year    int
	Season  int
	Episode int

	Resolution   Resolution
	Source       Source
	Language     Language
	Remux        bool
	Proper       bool
	DualLanguage bool
}

func (r *Result) IsMovie() bool {
	return r.MediaType == MediaTypeMovie
}

func (r *Result) IsTV() bool {
	return r.MediaType == MediaTypeTV
}

func (r *Result) VersionInfo() string {
	tokens := []string{}
	if r.Language != LangNA {
		tokens = append(tokens, r.Language.String())
	}
	if r.Resolution != ResNA {
		tokens = append(tokens, r.Resolution.String())
	}
	if r.DualLanguage {
		tokens = append(tokens, "DL")
	}
	if r.Source != SourceNA {
		tokens = append(tokens, r.Source.String())
	}
	if r.Remux {
		tokens = append(tokens, "Remux")
	}
	return strings.Join(tokens, ".")
}

func (r *Result) score() int {
	score := 0
	if r.Title != "" {
		score += 1
	}
	if r.MediaType != MediaTypeUnknown {
		score += 1
	}
	if r.Year != 0 {
		score += 1
	}
	if r.Season != 0 {
		score += 1
	}
	if r.Episode != 0 {
		score += 1
	}
	if r.Resolution != 0 {
		score += 1
	}
	if r.Source != SourceNA {
		score += 1
	}
	if r.Language != LangNA {
		score += 1
	}
	if r.Remux {
		score += 1
	}
	if r.Proper {
		score += 1
	}
	if r.DualLanguage {
		score += 1
	}
	return score
}

func (r *Result) mergeIn(other Result) {
	if other.Title != "" {
		r.Title = other.Title
	}
	if other.MediaType != MediaTypeUnknown {
		r.MediaType = other.MediaType
	}
	if other.Year != 0 {
		r.Year = other.Year
	}
	if other.Season != 0 {
		r.Season = other.Season
	}
	if other.Episode != 0 {
		r.Episode = other.Episode
	}
	if other.Resolution != 0 {
		r.Resolution = other.Resolution
	}
	if other.Source != SourceNA {
		r.Source = other.Source
	}
	if other.Language != LangNA {
		r.Language = other.Language
	}
	if other.Remux {
		r.Remux = other.Remux
	}
	if other.Proper {
		r.Proper = other.Proper
	}
	if other.DualLanguage {
		r.DualLanguage = other.DualLanguage
	}
}
