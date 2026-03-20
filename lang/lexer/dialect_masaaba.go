package lexer

import (
	"regexp"
)

// Country: Uganda
// Masaaba (Lumasaaba) is a Bantu language spoken by the Bamasaba people
// of eastern Uganda, primarily around Mount Elgon. It is a major language
// of the region with a rich cultural heritage linked to the land and tradition.
type Masaaba struct{}

var _ Dialect = (*Masaaba)(nil)

func (m Masaaba) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`masaaba:masaaba;`)
}

func (m Masaaba) Name() string {
	return "lumasaaba"
}

func (m Masaaba) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"leka":      Let,
		"nywevu":    Const,
		"singa":     KeywordIf,
		"ate":       KeywordElse,
		"nga":       KeywordWhile,
		"omulimu":   Function,
		"zaayo":     Return,
		"tandika":   Main,
		"lagako":    Print,
		"ulwakhiwo": Struct,
		"mazima":    True,
		"bulimba":   False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"float":     TypeFloat,
		"enamba":    TypeInt,
		"enamba64":  TypeInt64,
		"ebigambo":  TypeString,
		"ensobi":    TypeError,
	}
}

func (m Masaaba) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
