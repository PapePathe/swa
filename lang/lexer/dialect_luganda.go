package lexer

import (
	"regexp"
)

// Country: Uganda
// Luganda is a major Bantu language spoken primarily in the Buganda region
// of central Uganda. It is the most widely spoken second language in
// Uganda and is central to the country's administration, commerce, and media.
type Luganda struct{}

var _ Dialect = (*Luganda)(nil)

func (m Luganda) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`luganda:luganda;`)
}

func (m Luganda) Name() string {
	return "luganda"
}

func (m Luganda) Reserved() map[string]TokenKind {
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
		"desimal":   TypeFloat,
		"enamba":    TypeInt,
		"enamba64":  TypeInt64,
		"ebigambo":  TypeString,
		"ensobi":    TypeError,
	}
}

func (m Luganda) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
