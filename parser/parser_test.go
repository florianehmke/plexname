package parser_test

import (
	"testing"

	"github.com/florianehmke/plexname/parser"
)

type parserTest struct {
	toParse      string
	expectations parser.Result
}

var parserTests = []parserTest{
	{
		toParse:      "Some Title",
		expectations: parser.Result{},
	},
	{
		toParse:      "Some Title HDTV",
		expectations: parser.Result{Source: parser.HDTV},
	},
	{
		toParse:      "Some Title PDTV",
		expectations: parser.Result{Source: parser.PDTV},
	},
	{
		toParse:      "Some Title SDTV",
		expectations: parser.Result{Source: parser.SDTV},
	},
	{
		toParse:      "Some Title TVRip",
		expectations: parser.Result{Source: parser.TV},
	},
	{
		toParse:      "Some Title BD",
		expectations: parser.Result{Source: parser.BluRay},
	},
	{
		toParse:      "Some Title BR-Rip",
		expectations: parser.Result{Source: parser.BluRay},
	},
	{
		toParse:      "Some Title Blu-Ray",
		expectations: parser.Result{Source: parser.BluRay},
	},
	{
		toParse:      "Some Title DVD",
		expectations: parser.Result{Source: parser.DVD},
	},
	{
		toParse:      "Some Title.avi 720p",
		expectations: parser.Result{Resolution: parser.R720},
	},
	{
		toParse:      "Some Title.avi 720p webdl",
		expectations: parser.Result{Resolution: parser.R720, Source: parser.WEBDL},
	},
	{
		toParse:      "1080p web dl of Some Title",
		expectations: parser.Result{Resolution: parser.R1080, Source: parser.WEBDL},
	},
	{
		toParse:      "1080p.web-dl.of.A.Movie",
		expectations: parser.Result{Resolution: parser.R1080, Source: parser.WEBDL},
	},
	{
		toParse:      "Some Title repack",
		expectations: parser.Result{Proper: parser.True},
	},
	{
		toParse:      "Some.WEB-DL-HUNDUB.1080P",
		expectations: parser.Result{Resolution: parser.R1080, Source: parser.WEBDL, Language: parser.Hungarian},
	},
	{
		toParse:      "Some.Title.2012.Remux",
		expectations: parser.Result{Year: 2012, Remux: parser.True, Title: "some title"},
	},
	{
		toParse:      "Some.Title.S04E01",
		expectations: parser.Result{Season: 4, Episode1: 1, Title: "some title"},
	},
	{
		toParse:      "Some.Title.S01E02",
		expectations: parser.Result{Season: 1, Episode1: 2},
	},
	{
		toParse:      "Some.Title.WEB-DL",
		expectations: parser.Result{Source: parser.WEBDL},
	},
	{
		toParse:      "Some.Title.webrip",
		expectations: parser.Result{Source: parser.WEBRip},
	},
	{
		toParse:      "Some.Title.WEB-DL.DL",
		expectations: parser.Result{Source: parser.WEBDL, DualLanguage: parser.True},
	},
	{
		toParse:      "Some.Title.DL",
		expectations: parser.Result{DualLanguage: parser.True},
	},
	{
		toParse:      "Some.Title.E01E02",
		expectations: parser.Result{Episode1: 1, Episode2: 2},
	},
	{
		toParse: "Some.Title.S04E01.GERMAN.DL.1080p.BluRay/the200-1080p.mkv",
		expectations: parser.Result{
			Season:       4,
			Language:     parser.German,
			DualLanguage: parser.True,
			Resolution:   parser.R1080,
			Source:       parser.BluRay,
			Episode1:     1,
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range parserTests {
		t.Logf("Testing string: %s", test.toParse)
		got := parser.Parse(test.toParse, parser.Result{})
		compareResult(t, test.expectations, got)
	}
}

func compareResult(t *testing.T, expected parser.Result, got parser.Result) {
	if expected.Resolution != got.Resolution {
		t.Errorf("expected resolution=%s, got resolution=%s", expected.Resolution.String(), got.Resolution.String())
	}
	if expected.Source != got.Source {
		t.Errorf("expected source=%s, got source=%s", expected.Source.String(), got.Source.String())
	}
	if expected.Proper != got.Proper {
		t.Errorf("expected proper=%d, got proper=%d", expected.Proper, got.Proper)
	}
	if expected.Remux != got.Remux {
		t.Errorf("expected remux=%d, got remux=%d", expected.Remux, got.Remux)
	}
	if expected.Language != got.Language {
		t.Errorf("expected language=%d, got language=%d", expected.Language, got.Language)
	}
	if expected.Season != got.Season {
		t.Errorf("expected season=%d, got season=%d", expected.Season, got.Season)
	}
	if expected.Episode1 != got.Episode1 {
		t.Errorf("expected episode1=%d, got episode1=%d", expected.Episode1, got.Episode1)
	}
	if expected.Episode2 != got.Episode2 {
		t.Errorf("expected episode2=%d, got episode2=%d", expected.Episode2, got.Episode2)
	}
	if expected.Year != got.Year {
		t.Errorf("expected year=%d, got year=%d", expected.Year, got.Year)
	}
	if expected.Title != "" && expected.Title != got.Title {
		t.Errorf("expected title=%s, got title=%s", expected.Title, got.Title)
	}
	if expected.DualLanguage != got.DualLanguage {
		t.Errorf("expected dual-language=%d, got dual-language=%d", expected.DualLanguage, got.DualLanguage)
	}
}

func TestOverride(t *testing.T) {
	overrides := parser.Result{
		Title:        "Some Title",
		MediaType:    parser.MediaTypeTV,
		Year:         1999,
		Season:       2,
		Episode1:     15,
		Resolution:   parser.R2160,
		Source:       parser.BluRay,
		Language:     parser.German,
		Remux:        parser.True,
		Proper:       parser.False,
		DualLanguage: parser.True,
	}
	result := parser.Parse("1080p web dl of Some Title", overrides)
	if overrides != result {
		t.Errorf("expected overrides to have an effect")
	}
}
