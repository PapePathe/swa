package lexer

import (
	"regexp"
)

// Country: Uganda
// Rukiga (Chiga) is a Bantu language spoken by the Bakiga people in the Kigezi
// region of southwestern Uganda. It is part of the Runyakitara language cluster
// and is closely related to Runyankore.
type Rukiga struct{}

var _ Dialect = (*Rukiga)(nil)

func (m Rukiga) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`rukiga:rukiga;`)
}

func (m Rukiga) Name() string {
	return "rukiga"
}

func (m Rukiga) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"reka":       Let,
		"hamire":     Const,
		"ku":         KeywordIf,
		"ninga":      KeywordElse,
		"manya":      KeywordWhile,
		"omulimo":    Function,
		"garura":     Return,
		"tandika":    Main,
		"handika":    Print,
		"mutaratara": Struct,
		"mazima":     True,
		"bishuba":    False,
		"bool":       TypeBool,
		"bauto":      TypeByte,
		"float":      TypeFloat,
		"namba":      TypeInt,
		"namba64":    TypeInt64,
		"ebigambo":   TypeString,
		"ensobi":     TypeError,
	}
}

func (m Rukiga) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
