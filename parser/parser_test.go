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
		"/some/dir/Some Title.avi 720p webdl":  {Resolution: parser.R720, Source: parser.WEB},
		"/some/dir/1080p web dl of Some Title": {Resolution: parser.R1080, Source: parser.WEB},
		"/some/dir/1080p.web-dl.of.A.Movie":    {Resolution: parser.R1080, Source: parser.WEB},
		"/some/dir/Some Title repack":          {Proper: true},
		"/some/dir/Some.WEB.DL-HUNDUB.1080P":   {Resolution: parser.R1080, Source: parser.WEB, Language: parser.Hungarian},
		"/some/dir/Some.Title.2012.Remux":      {Year: 2012, Remux: true, Title: "some title"},
		"/some/dir/Some.Title.S01E02":          {Season: 1, Episode: 2},
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
		t.Errorf("expected proper=%t, got proper=%t", expected.Proper, got.Proper)
	}
	if expected.Remux != got.Remux {
		t.Errorf("expected remux=%t, got remux=%t", expected.Remux, got.Remux)
	}
	if expected.Language != got.Language {
		t.Errorf("expected language=%d, got language=%d", expected.Language, got.Language)
	}
	if expected.Season != got.Season {
		t.Errorf("expected season=%d, got season=%d", expected.Season, got.Season)
	}
	if expected.Episode != got.Episode {
		t.Errorf("expected episode=%d, got episode=%d", expected.Episode, got.Episode)
	}
	if expected.Year != got.Year {
		t.Errorf("expected year=%d, got year=%d", expected.Year, got.Year)
	}
	if expected.Title != "" && expected.Title != got.Title {
		t.Errorf("expected title=%s, got title=%s", expected.Title, got.Title)
	}
}

func TestOverride(t *testing.T) {
	overrides := parser.Result{
		Title:        "Some Title",
		MediaType:    parser.MediaTypeTV,
		Year:         1999,
		Season:       2,
		Episode:      15,
		Resolution:   parser.R2160,
		Source:       parser.BluRay,
		Language:     parser.German,
		Remux:        true,
		Proper:       true,
		DualLanguage: true,
	}
	result := parser.Parse("1080p web dl of Some Title", "/dev/null", overrides)
	if overrides != *result {
		t.Errorf("expected overrides to have an effect")
	}
}
