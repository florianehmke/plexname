package parser

import (
	"regexp"
	"strconv"
)

var (
	yearRegEx    = regexp.MustCompile(`(?P<year>(19|20)\d{2})`)
	seasonRegEx  = regexp.MustCompile(`s(?P<season>\d{1,2})`)
	episodeRegEx = regexp.MustCompile(`e(?P<episode>\d{2,4})`)

	fallbackRegExList = []*regexp.Regexp{
		// Show Title S01/1 - Title.mkv
		regexp.MustCompile(`.*s(?P<season>\d{1,2}).*/(?P<episode>\d{1,2}).+`),
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
	if match, ok := paramsMap["episode"]; ok {
		episode, _ := strconv.Atoi(match)
		result.Episode = episode
	}
	if match, ok := paramsMap["season"]; ok {
		season, _ := strconv.Atoi(match)
		result.Season = season
	}
	if match, ok := paramsMap["year"]; ok {
		year, _ := strconv.Atoi(match)
		result.Year = year
	}

	return result
}
