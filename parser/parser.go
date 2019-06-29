package parser

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type MediaType int

const (
	MediaTypeUnknown MediaType = iota
	MediaTypeMovie
	MediaTypeTV
)

var (
	yearRegEx    = regexp.MustCompile(`19|20\d{2}`)
	seasonRegEx  = regexp.MustCompile(`s\d{2}`)
	episodeRegEx = regexp.MustCompile(`e\d{2}`)

	remuxes = map[string]bool{
		"remux": true,
	}

	propers = map[string]bool{
		"repack": true,
		"rerip":  true,
		"proper": true,
	}
)

type Result struct {
	Title string

	MediaType MediaType

	Year    int
	Season  int
	Episode int

	Resolution Resolution
	Source     Source
	Remux      bool
	Proper     bool
}

func (r *Result) IsUnknownMediaType() bool {
	return r.MediaType == MediaTypeUnknown
}

func (r *Result) IsMovie() bool {
	return r.MediaType == MediaTypeMovie
}

func (r *Result) IsTV() bool {
	return r.MediaType == MediaTypeTV
}

func Parse(releaseName string, overrides Result) *Result {
	p := newParser(releaseName, overrides)
	p.parseTitle()
	p.parseYear()
	p.parseResolution()
	p.parseSource()
	p.parseRemux()
	p.parseProper()
	p.parseEpisode()
	p.parseSeason()
	p.parseMediaType()
	return p.result
}

type parser struct {
	overrides     *Result
	result        *Result
	releaseName   string
	cleanedName   string
	releaseTokens []string
}

func newParser(releaseName string, overrides Result) *parser {
	lr := strings.ToLower(releaseName)
	return &parser{
		overrides:     &overrides,
		releaseName:   lr,
		cleanedName:   clean(lr),
		releaseTokens: tokenize(lr),
		result:        &Result{},
	}
}

func (p *parser) parseTitle() {
	if p.overrides.Title != "" {
		p.result.Title = p.overrides.Title
		return
	}
	titleTokens := []string{}
	for _, t := range p.releaseTokens {
		if yearRegEx.MatchString(t) || seasonRegEx.MatchString(t) || episodeRegEx.MatchString(t) {
			break
		}
		titleTokens = append(titleTokens, t)
	}
	p.result.Title = strings.Join(titleTokens, " ")
}

func (p *parser) parseYear() {
	if p.overrides.Year != 0 {
		p.result.Year = p.overrides.Year
		return
	}
	for _, t := range p.releaseTokens {
		if yearRegEx.MatchString(t) {
			year, err := strconv.Atoi(t)
			if err == nil {
				p.result.Year = year
			}
		}
	}
}

func (p *parser) parseResolution() {
	if p.overrides.Resolution != ResNA {
		p.result.Resolution = p.overrides.Resolution
		return
	}
	for _, t := range p.releaseTokens {
		if res, ok := resMap[t]; ok {
			p.result.Resolution = res
		}
	}
	for k, res := range resMap {
		if strings.Contains(p.cleanedName, k) {
			p.result.Resolution = res
		}
	}
}

func (p *parser) parseSource() {
	if p.overrides.Source != SourceNA {
		p.result.Source = p.overrides.Source
		return
	}
	for _, t := range p.releaseTokens {
		if src, ok := srcMap[t]; ok {
			p.result.Source = src
		}
	}
	for k, src := range srcMap {
		// A source with less than 5 characters
		// produces too many false positives.
		if len(k) < 5 {
			continue
		}
		if strings.Contains(p.cleanedName, k) {
			p.result.Source = src
		}
	}
}

func (p *parser) parseRemux() {
	if p.overrides.Remux != false {
		p.result.Remux = p.overrides.Remux
		return
	}
	for _, t := range p.releaseTokens {
		if _, ok := remuxes[t]; ok {
			p.result.Remux = true
		}
	}
	for k := range remuxes {
		if strings.Contains(p.cleanedName, k) {
			p.result.Remux = true
		}
	}
}

func (p *parser) parseProper() {
	if p.overrides.Proper != false {
		p.result.Proper = p.overrides.Proper
		return
	}
	for _, t := range p.releaseTokens {
		if _, ok := propers[t]; ok {
			p.result.Proper = true
		}
	}
	for k := range propers {
		if strings.Contains(p.cleanedName, k) {
			p.result.Proper = true
		}
	}
}

func (p *parser) parseEpisode() {
	if p.overrides.Episode != 0 {
		p.result.Episode = p.overrides.Episode
		return
	}
	for _, t := range p.releaseTokens {
		e := episodeRegEx.FindString(t)
		if e != "" {
			ep, err := strconv.Atoi(e[1:])
			if err == nil {
				p.result.Episode = ep
			}
		}
	}
}

func (p *parser) parseSeason() {
	if p.overrides.Season != 0 {
		p.result.Season = p.overrides.Season
		return
	}
	for _, t := range p.releaseTokens {
		s := seasonRegEx.FindString(t)
		if s != "" {
			season, err := strconv.Atoi(s[1:])
			if err == nil {
				p.result.Season = season
			}
		}
	}
}

func (p *parser) parseMediaType() {
	if p.overrides.MediaType != MediaTypeUnknown {
		p.result.MediaType = p.overrides.MediaType
		return
	}
	if p.result.Episode != 0 && p.result.Season != 0 {
		p.result.MediaType = MediaTypeTV
	} else if p.result.Year != 0 {
		p.result.MediaType = MediaTypeMovie
	}
}

func isValidFileNameCharacter(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		return true
	}
	return false
}

func clean(s string) string {
	return strings.Map(
		func(r rune) rune {
			if !isValidFileNameCharacter(r) || unicode.IsSpace(r) {
				return -1
			}
			return r
		},
		s,
	)
}

func tokenize(s string) []string {
	t := strings.Map(
		func(r rune) rune {
			if !isValidFileNameCharacter(r) || unicode.IsSpace(r) {
				return rune(';')
			}
			return r
		},
		s,
	)
	return strings.Split(t, ";")
}
