package parser

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Parse(source, target string, overrides Result) *Result {
	srcPath, srcFile := filepath.Split(source)
	_, srcDir := filepath.Split(strings.TrimRight(srcPath, "/"))
	srcDirAndFile := srcDir + "/" + srcFile

	p := parser{
		sourcePath: newParseData(source),
		targetPath: newParseData(target),

		file:       newParseData(srcFile),
		dir:        newParseData(srcDir),
		dirAndFile: newParseData(srcDirAndFile),

		overrides: &overrides,
		result:    &Result{},
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
	p.parseMediaType()

	p.result.mergeIn(overrides)
	return p.result
}

type parser struct {
	sourcePath parseData
	targetPath parseData

	file       parseData
	dir        parseData
	dirAndFile parseData

	overrides *Result
	result    *Result
}

func (p *parser) dirOrFile() parseData {
	if p.file.length > p.dir.length {
		return p.file
	}
	if p.dir.length > p.file.length {
		return p.dir
	}
	return p.file
}

type parseData struct {
	name   string
	length int

	joined string
	tokens []string
}

func newParseData(s string) parseData {
	ls := strings.ToLower(s)
	return parseData{
		name:   ls,
		length: len(ls),
		joined: clean(ls),
		tokens: tokenize(ls),
	}
}

func (p *parser) parseTitle() {
	var titleTokens []string
	for _, t := range p.dirOrFile().tokens {
		if yearRegEx.MatchString(t) || seasonRegEx.MatchString(t) || episodeRegEx.MatchString(t) {
			break
		}
		titleTokens = append(titleTokens, t)
	}
	p.result.Title = strings.Join(titleTokens, " ")
}

func (p *parser) parseYear() {
	for _, t := range p.dirOrFile().tokens {
		if yearRegEx.MatchString(t) {
			year, err := strconv.Atoi(t)
			if err == nil {
				p.result.Year = year
			}
		}
	}
}

func (p *parser) parseResolution() {
	for _, t := range p.dirOrFile().tokens {
		if res, ok := resMap[t]; ok {
			p.result.Resolution = res
		}
	}
	for k, res := range resMap {
		if strings.Contains(p.dirOrFile().joined, k) {
			p.result.Resolution = res
		}
	}
}

func (p *parser) parseSource() {
	for _, t := range p.dirOrFile().tokens {
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
		if strings.Contains(p.dirOrFile().joined, k) {
			p.result.Source = src
		}
	}
}

func (p *parser) parseLanguage() {
	for _, t := range p.dirOrFile().tokens {
		if lang, ok := langMap[t]; ok {
			p.result.Language = lang
		}
	}
	for k, lang := range langMap {
		if strings.Contains(p.dirOrFile().joined, k) {
			p.result.Language = lang
		}
	}
}

func (p *parser) parseDualLanguage() {
	for _, t := range p.dirOrFile().tokens {
		if t == "dl" {
			count := strings.Count(p.dirOrFile().joined, "dl")
			webDL := strings.Contains(p.dirOrFile().joined, "webdl")
			if count > 1 || !webDL {
				p.result.DualLanguage = true
			}
		}
	}
}

func (p *parser) parseRemux() {
	for _, t := range p.dirOrFile().tokens {
		if t == "remux" {
			p.result.Remux = true
		}
	}
	if strings.Contains(p.dirOrFile().joined, "remux") {
		p.result.Remux = true
	}
}

func (p *parser) parseProper() {
	propers := map[string]bool{
		"repack": true,
		"rerip":  true,
		"proper": true,
	}

	for _, t := range p.dirOrFile().tokens {
		if _, ok := propers[t]; ok {
			p.result.Proper = true
		}
	}
	for k := range propers {
		if strings.Contains(p.dirOrFile().joined, k) {
			p.result.Proper = true
		}
	}
}

func (p *parser) parseSeasonAndEpisode() {
	for _, t := range p.dirOrFile().tokens {
		r := populateResultFromRxpList([]*regexp.Regexp{seasonRegEx, episodeRegEx}, t)
		p.result.mergeIn(r)
	}
	if p.result.Episode == 0 || p.result.Season == 0 {
		r := getBestResultFromRxpList(fallbackRegExList, p.dirAndFile.name)
		if r.score() > 0 {
			p.result.Season = r.Season
			p.result.Episode = r.Episode
		}
	}
}

func (p *parser) parseMediaType() {
	for _, tokens := range [][]string{
		p.targetPath.tokens,
		p.sourcePath.tokens,
	} {
		for _, s := range tokens {
			if mt, ok := mediaTypes[s]; ok {
				p.result.MediaType = mt
				return
			}
		}
	}
	if p.result.Episode != 0 || p.result.Season != 0 {
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
