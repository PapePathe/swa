package lexer

import (
	"regexp"
)

// Country: Mali
// Dogon is a cluster of languages spoken by the Dogon people of Mali. It is known
// for its unique linguistic features and its central role in the complex
// cosmological and social structures of the Dogon plateau.
type Dogon struct{}

var _ Dialect = (*Dogon)(nil)

func (m Dogon) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dogon:dogon;`)
}

func (m Dogon) Name() string {
	return "dogoro"
}

func (m Dogon) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"liye":    Let,
		"tiga":    Const,
		"ne":      KeywordIf,
		"ne-ma":   KeywordElse,
		"teme":    KeywordWhile,
		"kaadu":   Function,
		"yali":    Return,
		"dama":    Main,
		"tege":    Print,
		"caakaar": Struct,
		"ada":     True,
		"go-ba":   False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"desimal": TypeFloat,
		"lim":     TypeInt,
		"lim64":   TypeInt64,
		"sembe":   TypeString,
		"kile":    TypeError,
	}
}

func (m Dogon) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
