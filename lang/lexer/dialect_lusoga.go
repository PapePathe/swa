package lexer

import (
	"regexp"
)

// Country: Uganda
// Lusoga is a Bantu language spoken by the Basoga people in the Busoga region
// of eastern Uganda. It is one of the major languages of Uganda with a
// significant number of native speakers and a rich linguistic history.
type Lusoga struct{}

var _ Dialect = (*Lusoga)(nil)

func (m Lusoga) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`lusoga:lusoga;`)
}

func (m Lusoga) Name() string {
	return "lusoga"
}

func (m Lusoga) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"leka":      Let,
		"nywevu":    Const,
		"singa":     KeywordIf,
		"oba":       KeywordElse,
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

func (m Lusoga) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
