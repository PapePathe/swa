package lexer

import (
	"regexp"
)

// Country: Nigeria
// Kanuri is a major Nilo-Saharan language spoken primarily in the Lake Chad Basin, particularly
// in Borno and Yobe states of North-Eastern Nigeria. It has a long history as the language of
// the Kanem-Bornu Empire and remains a vital cultural identifier in the region.
type Kanuri struct{}

var _ Dialect = (*Kanuri)(nil)

func (m Kanuri) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`kanuri:kanuri;`)
}

func (m Kanuri) Name() string {
	return "kanuri"
}

func (m Kanuri) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"che":     Let,
		"tabat":   Const,
		"ya":      KeywordIf,
		"gade":    KeywordElse,
		"waati":   KeywordWhile,
		"gyara":   Function,
		"walta":   Return,
		"piilgu":  Main,
		"andika":  Print,
		"laben":   Struct,
		"kal":     True,
		"tangari": False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"decimal": TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"hantum":  TypeString,
		"daray":   TypeError,
	}
}

func (m Kanuri) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
