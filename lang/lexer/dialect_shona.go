package lexer

import (
	"regexp"
)

// Country: Zimbabwe
// Shona (chiShona) is a major Bantu language spoken by the Shona people of
// Zimbabwe. It is the most widely spoken native language in the country and
// has a significant literary tradition dating back centuries.
type Shona struct{}

var _ Dialect = (*Shona)(nil)

func (m Shona) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`shona:shona;`)
}

func (m Shona) Name() string {
	return "shona"
}

func (m Shona) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"nga":                 Let,
		"chisingashanduki":    Const,
		"kana":                KeywordIf,
		"kana_zvisina_kudaro": KeywordElse,
		"dzokorora":           KeywordWhile,
		"basa":                Function,
		"dzosa":               Return,
		"kutanga":             Main,
		"nyora":               Print,
		"chivakwa":            Struct,
		"chokwadi":            True,
		"nhema":               False,
		"bool":                TypeBool,
		"byte":                TypeByte,
		"float":               TypeFloat,
		"nhamba":              TypeInt,
		"nhamba64":            TypeInt64,
		"mazwi":               TypeString,
		"mhosva":              TypeError,
	}
}

func (m Shona) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
