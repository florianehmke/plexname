package parser_test

import (
	"testing"

	"github.com/florianehmke/plexname/parser"
)

var (
	qualityParserTests = map[string]parser.Result{
		"/some/dir/Some Title":                 {},
		"/some/dir/Some Title HDTV":            {Source: parser.HDTV},
		"/some/dir/Some Title PDTV":            {Source: parser.PDTV},
		"/some/dir/Some Title SDTV":            {Source: parser.SDTV},
		"/some/dir/Some Title TVRip":           {Source: parser.TV},
		"/some/dir/Some Title BD":              {Source: parser.BluRay},
		"/some/dir/Some Title BR-Rip":          {Source: parser.BluRay},
		"/some/dir/Some Title Blu-Ray":         {Source: parser.BluRay},
		"/some/dir/Some Title DVD":             {Source: parser.DVD},
		"/some/dir/Some Title.avi 720p":        {Resolution: parser.R720},
		"/some/dir/Some Title.avi 720p webdl":  {Resolution: parser.R720, Source: parser.WEBDL},
		"/some/dir/1080p web dl of Some Title": {Resolution: parser.R1080, Source: parser.WEBDL},
		"/some/dir/1080p.web-dl.of.A.Movie":    {Resolution: parser.R1080, Source: parser.WEBDL},
		"/some/dir/Some Title repack":          {Proper: parser.True},
		"/some/dir/Some.WEB-DL-HUNDUB.1080P":   {Resolution: parser.R1080, Source: parser.WEBDL, Language: parser.Hungarian},
		"/some/dir/Some.Title.2012.Remux":      {Year: 2012, Remux: parser.True, Title: "some title"},
		"/some/dir/Some.Title.S01E02":          {Season: 1, Episode1: 2},
		"/some/dir/Some.Title.WEB-DL":          {Source: parser.WEBDL},
		"/some/dir/Some.Title.webrip":          {Source: parser.WEBRip},
		"/some/dir/Some.Title.WEB-DL.DL":       {Source: parser.WEBDL, DualLanguage: parser.True},
		"/some/dir/Some.Title.DL":              {DualLanguage: parser.True},
		"/some/dir/Some.Title.E01E02":          {Episode1: 1, Episode2: 2},
	}
)

func TestParse(t *testing.T) {
	for title, expected := range qualityParserTests {
		t.Logf("Testing title: %s", title)
		got := parser.Parse(title, "/dev/null", parser.Result{})
		compareResult(t, &expected, got)
	}
}

func compareResult(t *testing.T, expected *parser.Result, got *parser.Result) {
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
	result := parser.Parse("1080p web dl of Some Title", "/dev/null", overrides)
	if overrides != *result {
		t.Errorf("expected overrides to have an effect")
	}
}
