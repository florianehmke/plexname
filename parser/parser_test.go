package parser_test

import (
	"testing"

	"github.com/florianehmke/plexname/parser"
)

var (
	qualityParserTests = map[string]parser.Result{
		"Some Title":                 parser.Result{},
		"Some Title HDTV":            parser.Result{Source: parser.HDTV},
		"Some Title PDTV":            parser.Result{Source: parser.PDTV},
		"Some Title SDTV":            parser.Result{Source: parser.SDTV},
		"Some Title TVRip":           parser.Result{Source: parser.TV},
		"Some Title BD":              parser.Result{Source: parser.BluRay},
		"Some Title BR-Rip":          parser.Result{Source: parser.BluRay},
		"Some Title Blu-Ray":         parser.Result{Source: parser.BluRay},
		"Some Title DVD":             parser.Result{Source: parser.DVD},
		"Some Title.avi 720p":        parser.Result{Resolution: parser.R720},
		"Some Title.avi 720p webdl":  parser.Result{Resolution: parser.R720, Source: parser.WEB},
		"1080p web dl of Some Title": parser.Result{Resolution: parser.R1080, Source: parser.WEB},
		"1080p.web-dl.of.A.Movie":    parser.Result{Resolution: parser.R1080, Source: parser.WEB},
		"Some Title repack":          parser.Result{Proper: true},
		"Some.WEB.DL-HUNDUB.1080P":   parser.Result{Resolution: parser.R1080, Source: parser.WEB},
		"Some.Title.2012.Remux":      parser.Result{Year: 2012, Remux: true, Title: "some title"},
		"Some.Title.S01E02":          parser.Result{Season: 1, Episode: 2},
	}
)

func TestParse(t *testing.T) {
	for title, expected := range qualityParserTests {
		t.Logf("Testing title: %s", title)
		got := parser.Parse(title)
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
