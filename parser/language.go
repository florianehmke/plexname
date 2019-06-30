package parser

import (
	"fmt"
	"strings"
)

// Language is the internal representation of a language.
type Language int

// All known languages.
const (
	LangNA Language = iota
	English
	French
	Spanish
	German
	Italian
	Danish
	Dutch
	Japanese
	Cantonese
	Mandarin
	Russian
	Polish
	Vietnamese
	Swedish
	Norwegian
	Finnish
	Turkish
	Portuguese
	Flemish
	Greek
	Korean
	Hungarian
)

var (
	langMap = map[string]Language{
		"english":    English,
		"fr":         French,
		"vostfr":     French,
		"french":     French,
		"spanish":    Spanish,
		"videomann":  German,
		"german":     German,
		"ita":        Italian,
		"italian":    Italian,
		"danish":     Danish,
		"nl":         Dutch,
		"dutch":      Dutch,
		"japanese":   Japanese,
		"cantonese":  Cantonese,
		"mandarin":   Mandarin,
		"rus":        Russian,
		"brus":       Russian,
		"russian":    Russian,
		"polish":     Polish,
		"vietnamese": Vietnamese,
		"swedish":    Swedish,
		"norwegian":  Norwegian,
		"finnish":    Finnish,
		"turkish":    Turkish,
		"portuguese": Portuguese,
		"flemish":    Flemish,
		"greek":      Greek,
		"korean":     Korean,
		"hun":        Hungarian,
		"hundub":     Hungarian,
		"hungarian":  Hungarian,
	}

	langNames = map[Language]string{
		LangNA:     "--",
		English:    "English",
		French:     "French",
		Spanish:    "Spanish",
		German:     "German",
		Italian:    "Italian",
		Danish:     "Danish",
		Dutch:      "Dutch",
		Japanese:   "Japanese",
		Cantonese:  "Cantonese",
		Mandarin:   "Mandarin",
		Russian:    "Russian",
		Polish:     "Polish",
		Vietnamese: "Vietnamese",
		Swedish:    "Swedish",
		Norwegian:  "Norwegian",
		Finnish:    "Finnish",
		Turkish:    "Turkish",
		Portuguese: "Portuguese",
		Flemish:    "Flemish",
		Greek:      "Greek",
		Korean:     "Korean",
		Hungarian:  "Hungarian",
	}
)

// ParseLanguage parses the given string to a language.
func ParseLanguage(lang string) (Language, error) {
	if l, ok := langMap[strings.ToLower(lang)]; ok {
		return l, nil
	}
	if lang == "" {
		return LangNA, nil
	}
	return LangNA, fmt.Errorf("unknown language: %s", lang)
}

// String returns the string representation of l.
func (l Language) String() string {
	return langNames[l]
}

// LanguageNames returns a map containing all string representations.
func LanguageNames() map[Language]string {
	return langNames
}
