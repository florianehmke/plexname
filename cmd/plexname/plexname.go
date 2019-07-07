package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/florianehmke/plexname/config"
	"github.com/florianehmke/plexname/fs"
	"github.com/florianehmke/plexname/log"
	"github.com/florianehmke/plexname/prompt"
	"github.com/florianehmke/plexname/renamer"
	"github.com/florianehmke/plexname/search"
	"github.com/florianehmke/plexname/tmdb"
	"github.com/florianehmke/plexname/tvdb"
)

func main() {
	flag.Usage = usage

	arguments := renamer.GetArgsFromFlags()
	r := renamer.New(
		arguments,
		search.NewSearcher(
			tmdb.NewClient(tmdb.BaseURL, config.GetToken("tmdb")),
			tvdb.NewClient(tvdb.BaseURL, config.GetToken("tvdb")),
			prompt.NewPrompter(),
		),
		fs.NewFileSystem(arguments.DryRun),
	)

	if err := r.Run(); err != nil {
		log.Error(fmt.Sprintf("renaming failed: %v", err))
		os.Exit(1)
	}

	log.Info("Yay, done!")
	os.Exit(0)
}

func usage() {
	fmt.Println("plexname")
	fmt.Println("  Rename your media files and folders for the Plex Media Server.")
	fmt.Println()
	fmt.Println("Usage: ")
	fmt.Println("  plexname [option]... source-dir [target-dir]")
	fmt.Println("  plexname [option]... file")
	fmt.Println("")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  plexname -extensions=mkv,mp4 -lang english -remux downloads movies")
}
