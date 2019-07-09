package parser

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Parse(s string, overrides Result) Result {
	p := parser{
		parseData: newParseData(s),
		overrides: overrides,
		result:    Result{},
	}

	p.parseTitle()
	p.parseYear()
	p.parseResolution()
	p.parseSource()
	p.parseLanguage()
	p.parseDualLanguage()
	p.parseRemux()
	p.parseProper()
	p.parseSeasonAndEpisode()
	p.parseSpecial()

	p.doPlausibilityCheck()
	p.setMediaType()

	p.result.mergeIn(overrides)
	return p.result
}

type parser struct {
	parseData parseData

	overrides Result
	result    Result
}

type parseData struct {
	toParse string
	length  int

	joined string
	tokens []string
}

func newParseData(s string) parseData {
	ls := strings.ToLower(s)
	return parseData{
		toParse: ls,
		length:  len(ls),
		joined:  clean(ls),
		tokens:  tokenize(ls),
	}
}

func (p *parser) parseMediaType() {
	for _, s := range p.parseData.tokens {
		if mt, ok := mediaTypes[s]; ok {
			p.result.MediaType = mt
			return
		}
	}
	if episodeRegEx.sub.MatchString(p.parseData.joined) && seasonRegEx.sub.MatchString(p.parseData.joined) {
		p.result.MediaType = MediaTypeTV
	} else {
		for _, r := range tvAlternativeRegExList {
			if r.MatchString(p.parseData.toParse) {
				p.result.MediaType = MediaTypeTV
				return
			}
		}
		p.result.MediaType = MediaTypeMovie
	}
}

func (p *parser) parseTitle() {
	var titleTokens []string
	for _, t := range p.parseData.tokens {
		if yearRegEx.full.MatchString(t) || singleEpisode.matchFull(t) || dualEpisode.matchFull(t) {
			break
		}
		titleTokens = append(titleTokens, t)
	}
	p.result.Title = strings.Join(titleTokens, " ")
}

func (p *parser) parseYear() {
	for _, t := range p.parseData.tokens {
		if yearRegEx.full.MatchString(t) {
			year, err := strconv.Atoi(t)
			if err == nil {
				p.result.Year = year
			}
		}
	}
}

func (p *parser) parseResolution() {
	for _, t := range p.parseData.tokens {
		if res, ok := resMap[t]; ok {
			p.result.Resolution = res
		}
	}
	for k, res := range resMap {
		if strings.Contains(p.parseData.joined, k) {
			p.result.Resolution = res
		}
	}
}

func (p *parser) parseSource() {
	for _, t := range p.parseData.tokens {
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
		if strings.Contains(p.parseData.joined, k) {
			p.result.Source = src
		}
	}
}

func (p *parser) parseLanguage() {
	for _, t := range p.parseData.tokens {
		if lang, ok := langMap[t]; ok {
			p.result.Language = lang
		}
	}
	for k, lang := range langMap {
		if len(k) < 5 {
			continue
		}
		if strings.Contains(p.parseData.joined, k) {
			p.result.Language = lang
		}
	}
}

func (p *parser) parseDualLanguage() {
	for _, t := range p.parseData.tokens {
		if t == "dl" {
			count := strings.Count(p.parseData.joined, "dl")
			webDL := strings.Contains(p.parseData.joined, "webdl")
			if (!webDL && count == 1) || count > 1 {
				p.result.DualLanguage = True
			}
		}
	}
}

func (p *parser) parseRemux() {
	for _, t := range p.parseData.tokens {
		if t == "remux" {
			p.result.Remux = True
		}
	}
	if strings.Contains(p.parseData.joined, "remux") {
		p.result.Remux = True
	}
}

func (p *parser) parseProper() {
	propers := map[string]bool{
		"repack": true,
		"rerip":  true,
		"proper": true,
	}

	for _, t := range p.parseData.tokens {
		if _, ok := propers[t]; ok {
			p.result.Proper = True
		}
	}
	for k := range propers {
		if strings.Contains(p.parseData.joined, k) {
			p.result.Proper = True
		}
	}
}

func (p *parser) parseSeasonAndEpisode() {
	for _, t := range p.parseData.tokens {
		var r Result
		if dualEpisode.episode.sub.MatchString(t) {
			r = populateResultFromRxpList([]*regexp.Regexp{
				dualEpisode.episode.full,
				dualEpisode.complete.full,
			}, t)
		} else {
			r = populateResultFromRxpList([]*regexp.Regexp{
				singleEpisode.episode.full,
				singleEpisode.season.full,
				singleEpisode.complete.full,
			}, t)
		}
		p.result.mergeIn(r)
	}
	if p.result.Episode1 == 0 || p.result.Season == 0 {
		r := getBestResultFromRxpList(tvAlternativeRegExList, p.parseData.toParse)
		if r.score() > 0 {
			p.result.Season = r.Season
			p.result.Episode1 = r.Episode1
			p.result.Episode2 = r.Episode2
		}
	}
}

func (p *parser) parseSpecial() {
	specials := map[string]bool{
		"special":  true,
		"specials": true,
		"extra":    true,
		"extras":   true,
	}

	for _, t := range p.parseData.tokens {
		if _, ok := specials[t]; ok {
			p.result.Special = True
		}
	}
	if strings.Contains(p.parseData.joined, "s00") {
		p.result.Special = True
	}
}

func (p *parser) doPlausibilityCheck() {
	// If a season (>0) is present, it can't be a special.
	if p.result.Season > 0 && p.result.Special == True {
		p.result.Special = False
	}
}

func (p *parser) setMediaType() {
	if (p.result.Episode1 > 0 && p.result.Season > 0) || p.result.Special == True {
		p.result.MediaType = MediaTypeTV
	} else {
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
