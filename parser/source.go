package parser

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
	WEB
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
		WEB:      "WEB",
		DSR:      "DSR",
	}

	srcScores = map[Source]int{
		BluRay: 100,
		DVD:    20,
		TV:     10,
		HDTV:   30,
		SDTV:   10,
		PDTV:   10,
		WEB:    50,
		DSR:    10,
	}

	srcMap = map[string]Source{
		"bluray":   BluRay,
		"hddvd":    BluRay,
		"bd":       BluRay,
		"bdrip":    BluRay,
		"brrip":    BluRay,
		"webdl":    WEB,
		"webrip":   WEB,
		"ituneshd": WEB,
		"webhd":    WEB,
		"webx264":  WEB,
		"webx265":  WEB,
		"webh264":  WEB,
		"webh265":  WEB,
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

// String returns the string representation of r.
func (s Source) String() string {
	return srcNames[s]
}
