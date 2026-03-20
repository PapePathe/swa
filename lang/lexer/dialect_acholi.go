package lexer

import (
	"regexp"
)

// Country: Uganda
// Acholi is a Nilotic language spoken by the Acholi people in Northern
// Uganda and parts of South Sudan. it is a member of the Western Nilotic
// group and has a strong tradition of oral poetry and song.
type Acholi struct{}

var _ Dialect = (*Acholi)(nil)

func (m Acholi) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`acholi:acholi;`)
}

func (m Acholi) Name() string {
	return "acholi"
}

func (m Acholi) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"wek":         Let,
		"ma-pe-lokke": Const,
		"ka":          KeywordIf,
		"ka-pe":       KeywordElse,
		"ka-podi":     KeywordWhile,
		"tic":         Function,
		"dwok-cen":    Return,
		"cak-ki":      Main,
		"go-piny":     Print,
		"yub":         Struct,
		"ada":         True,
		"go-ba":       False,
		"bool":        TypeBool,
		"byte":        TypeByte,
		"desimal":     TypeFloat,
		"namba":       TypeInt,
		"namba64":     TypeInt64,
		"lok-anyen":   TypeString,
		"bal":         TypeError,
	}
}

func (m Acholi) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
