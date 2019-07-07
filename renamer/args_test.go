package renamer_test

import (
	"os"
	"testing"

	"github.com/florianehmke/plexname/parser"
	"github.com/florianehmke/plexname/renamer"
)

func TestGetArgsFromFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"plexname",
		"-dry",
		"-dl", "true",
		"-episode1", "5",
		"-episode2", "6",
		"-extensions", "mkv,avi",
		"-lang", "german",
		"-media-type", "tv",
		"-proper", "true",
		"-remux", "true",
		"-resolution", "720p",
		"-season", "3",
		"-source", "bluray",
		"-title", "Some Title",
		"-year", "1999",
		"-only-dir",
		"-only-file",
		"some/path",
		"some/other/path",
	}
	args := renamer.GetArgsFromFlags()
	if args.SourcePath != "some/path" {
		t.Error("expected different source path")
	}
	if args.TargetPath != "some/other/path" {
		t.Error("expected different target path")
	}
	if args.DryRun != true {
		t.Error("expected different dryRun flag value")
	}
	if len(args.Extensions) != 2 {
		t.Error("expected 2 extensions")
	}
	if !args.OnlyDir {
		t.Error("expected -only-dir to have an effect")
	}
	if !args.OnlyFile {
		t.Error("expected -only-file to have an effect")
	}
	expectedOverrides := parser.Result{
		Title:        "Some Title",
		DualLanguage: parser.True,
		Episode1:     5,
		Episode2:     6,
		Language:     parser.German,
		MediaType:    parser.MediaTypeTV,
		Proper:       parser.True,
		Remux:        parser.True,
		Resolution:   parser.R720,
		Season:       3,
		Source:       parser.BluRay,
		Year:         1999,
	}
	if args.Overrides != expectedOverrides {
		t.Error("expected different overrides")
	}
}
