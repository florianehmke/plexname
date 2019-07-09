package parser

import (
	"fmt"
	"strings"
)

// Source is the internal representation of a media source such as DVD.
type Source int

// All known media sources.
const (
	SourceNA Source = iota
	BluRay
	DVD
	TV
	HDTV
	SDTV
	PDTV
	WEBRip
	WEBDL
	DSR
)

// Sources mapped to their string representations.
var (
	srcNames = map[Source]string{
		SourceNA: "--",
		BluRay:   "Blu-ray",
		DVD:      "DVD",
		TV:       "TV",
		HDTV:     "HDTV",
		SDTV:     "SDTV",
		PDTV:     "PDTV",
		WEBRip:   "WEB-Rip",
		WEBDL:    "WEB-DL",
		DSR:      "DSR",
	}

	srcMap = map[string]Source{
		"bluray":   BluRay,
		"hddvd":    BluRay,
		"bd":       BluRay,
		"bdrip":    BluRay,
		"brrip":    BluRay,
		"webdl":    WEBDL,
		"ituneshd": WEBDL,
		"amzn":     WEBDL,
		"amazonhd": WEBDL,
		"webrip":   WEBRip,
		"webhd":    WEBRip,
		"webx264":  WEBRip,
		"webx265":  WEBRip,
		"webh264":  WEBRip,
		"webh265":  WEBRip,
		"tvrip":    TV,
		"hdtv":     HDTV,
		"pdtv":     PDTV,
		"sdtv":     SDTV,
		"wsdsr":    DSR,
		"dsr":      DSR,
		"dvd":      DVD,
		"dvdrip":   DVD,
		"ntsc":     DVD,
		"pal":      DVD,
		"xvidvd":   DVD,
	}
)

// ParseSource parses the given source to a sour e.
func ParseSource(source string) (Source, error) {
	if s, ok := srcMap[strings.ToLower(source)]; ok {
		return s, nil
	}
	if source == "" {
		return SourceNA, nil
	}
	return SourceNA, fmt.Errorf("unknown source: %s", source)
}

// String returns the string representation of r.
func (s Source) String() string {
	return srcNames[s]
}
