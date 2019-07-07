package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type regex struct {
	full *regexp.Regexp
	sub  *regexp.Regexp
}

type tvrxp struct {
	season   regex // ..S01..
	episode  regex // ..E01..
	complete regex // ..S01E02..
}

func mustCompile(str string) regex {
	return regex{
		full: regexp.MustCompile(fmt.Sprintf("^%s$", str)),
		sub:  regexp.MustCompile(str),
	}
}

var (
	seasonPattern   = `s(?P<season>\d{1,2})`
	episode1Pattern = `e(?P<episode1>\d{2,4})`
	episode2Pattern = `e(?P<episode2>\d{2,4})`

	yearRegEx    = mustCompile(`(?P<year>(19|20)\d{2})`)
	seasonRegEx  = mustCompile(seasonPattern)
	episodeRegEx = mustCompile(episode1Pattern)

	singleEpisode = tvrxp{
		season:   mustCompile(seasonPattern),
		episode:  mustCompile(episode1Pattern),
		complete: mustCompile(seasonPattern + episode1Pattern),
	}

	dualEpisode = tvrxp{
		season:   mustCompile(seasonPattern),
		episode:  mustCompile(episode1Pattern + episode2Pattern),
		complete: mustCompile(seasonPattern + episode1Pattern + episode2Pattern),
	}

	tvAlternativeRegExList = []*regexp.Regexp{
		// Show Title S01/1 - Title.mkv
		regexp.MustCompile(`.*s(?P<season>\d{1,2}).*/(?P<episode1>\d{1,4}).+`),
	}
)

func populateResultFromRxpList(rxps []*regexp.Regexp, s string) Result {
	result := Result{}
	for _, rxp := range rxps {
		r := getResultFromRegEx(rxp, s)
		result.mergeIn(r)
	}
	return result
}

func getBestResultFromRxpList(rxps []*regexp.Regexp, s string) Result {
	var results []Result
	for _, rxp := range rxps {
		result := getResultFromRegEx(rxp, s)
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

func getResultFromRegEx(rxp *regexp.Regexp, s string) Result {
	match := rxp.FindStringSubmatch(s)
	paramsMap := map[string]string{}
	for i, name := range rxp.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	result := Result{}
	if match, ok := paramsMap["season"]; ok {
		season, _ := strconv.Atoi(match)
		result.Season = season
	}
	if match, ok := paramsMap["year"]; ok {
		year, _ := strconv.Atoi(match)
		result.Year = year
	}
	if match, ok := paramsMap["episode1"]; ok {
		episode, _ := strconv.Atoi(match)
		result.Episode1 = episode
	}
	if match, ok := paramsMap["episode2"]; ok {
		episode, _ := strconv.Atoi(match)
		result.Episode2 = episode
	}
	return result
}
