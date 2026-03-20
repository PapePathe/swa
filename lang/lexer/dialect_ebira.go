package lexer

import (
	"regexp"
)

// Country: Nigeria
// Ebira is a Nupoid language spoken by the Ebira people, primarily in Kogi State around
// the Niger-Benue confluence. It is a major language of North-Central Nigeria with a
// strong cultural presence in the region.
type Ebira struct{}

var _ Dialect = (*Ebira)(nil)

func (m Ebira) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`ebira:ebira;`)
}

func (m Ebira) Name() string {
	return "ebira"
}

func (m Ebira) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"gi":      Let,
		"tabat":   Const,
		"if":      KeywordIf,
		"else":    KeywordElse,
		"while":   KeywordWhile,
		"tuma":    Function,
		"labi":    Return,
		"start":   Main,
		"sabbi":   Print,
		"rny":     Struct,
		"gaskiya": True,
		"karya":   False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"decimal": TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"mo":      TypeString,
		"obo":     TypeError,
	}
}

func (m Ebira) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
