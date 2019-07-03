package parser

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	yearRegEx    = regexp.MustCompile(`19|20\d{2}`)
	seasonRegEx  = regexp.MustCompile(`s\d{2}`)
	episodeRegEx = regexp.MustCompile(`e\d{2}`)

	fallbackRegExList = []*regexp.Regexp{
		// Show Title S01/1 - Title.mkv
		regexp.MustCompile(`.*s(?P<season>\d{1,2}).*/(?P<episode>\d{1,2}).+`),
	}

	remuxes = map[string]bool{
		"remux": true,
	}

	dualLangs = map[string]bool{
		"dl": true,
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

	Resolution   Resolution
	Source       Source
	Language     Language
	Remux        bool
	Proper       bool
	DualLanguage bool
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

func (r *Result) IncompleteTVResult() bool {
	return r.MediaType == MediaTypeTV && r.Season == 0 && r.Episode == 0
}

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

	p.parseMediaType()
	p.parseTitle()
	p.parseYear()
	p.parseResolution()
	p.parseSource()
	p.parseLanguage()
	p.parseDualLanguage()
	p.parseRemux()
	p.parseProper()
	p.parseSeasonAndEpisode()
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

func (p *parser) parseMediaType() {
	if p.overrides.MediaType != MediaTypeUnknown {
		p.result.MediaType = p.overrides.MediaType
		return
	}
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
	if seasonRegEx.MatchString(p.file.name) || seasonRegEx.MatchString(p.dir.name) {
		p.result.MediaType = MediaTypeTV
		return
	}
	if episodeRegEx.MatchString(p.file.name) || episodeRegEx.MatchString(p.dir.name) {
		p.result.MediaType = MediaTypeTV
		return
	}
	if yearRegEx.MatchString(p.file.name) || yearRegEx.MatchString(p.dir.name) {
		p.result.MediaType = MediaTypeMovie
		return
	}
}

func (p *parser) parseTitle() {
	if p.overrides.Title != "" {
		p.result.Title = p.overrides.Title
		return
	}
	titleTokens := []string{}
	for _, t := range p.dirOrFile().tokens {
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
	if p.overrides.Resolution != ResNA {
		p.result.Resolution = p.overrides.Resolution
		return
	}
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
	if p.overrides.Source != SourceNA {
		p.result.Source = p.overrides.Source
		return
	}
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
	if p.overrides.Language != LangNA {
		p.result.Language = p.overrides.Language
		return
	}
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
	if p.overrides.DualLanguage != false {
		p.result.DualLanguage = p.overrides.DualLanguage
		return
	}
	for _, t := range p.dirOrFile().tokens {
		if _, ok := dualLangs[t]; ok {
			p.result.DualLanguage = true
		}
	}
}

func (p *parser) parseRemux() {
	if p.overrides.Remux != false {
		p.result.Remux = p.overrides.Remux
		return
	}
	for _, t := range p.dirOrFile().tokens {
		if _, ok := remuxes[t]; ok {
			p.result.Remux = true
		}
	}
	for k := range remuxes {
		if strings.Contains(p.dirOrFile().joined, k) {
			p.result.Remux = true
		}
	}
}

func (p *parser) parseProper() {
	if p.overrides.Proper != false {
		p.result.Proper = p.overrides.Proper
		return
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
	if p.overrides.Episode != 0 {
		p.result.Episode = p.overrides.Episode
		return
	}
	for _, t := range p.dirOrFile().tokens {
		e := episodeRegEx.FindString(t)
		if e != "" {
			ep, err := strconv.Atoi(e[1:])
			if err == nil {
				p.result.Episode = ep
			}
		}
		s := seasonRegEx.FindString(t)
		if s != "" {
			season, err := strconv.Atoi(s[1:])
			if err == nil {
				p.result.Season = season
			}
		}
	}

	if p.result.Episode == 0 || p.result.Season == 0 {
		r := getResultFromFallbackRegEx(p.dirAndFile.name)
		if r.score() > 0 {
			p.result.Season = r.Season
			p.result.Episode = r.Episode
		}
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

func getResultFromFallbackRegEx(s string) Result {
	var results []Result
	for _, r := range fallbackRegExList {
		match := r.FindStringSubmatch(s)
		paramsMap := map[string]string{}
		for i, name := range r.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}

		result := Result{}
		if match, ok := paramsMap["episode"]; ok {
			episode, _ := strconv.Atoi(match)
			result.Episode = episode
		}
		if match, ok := paramsMap["season"]; ok {
			season, _ := strconv.Atoi(match)
			result.Season = season
		}

		results = append(results, result)
	}

	result, score := Result{}, 0
	for _, r := range results {
		if r.score() > score {
			result, score = r, r.score()
		}
	}
	return result
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
