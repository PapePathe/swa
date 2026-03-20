package lexer

import (
	"regexp"
)

// Country: Madagascar
// Malagasy is an Austronesian language and the national language of
// Madagascar. It is unique in Africa for its Austronesian roots, mixed
// with Bantu, Arabic, and European influences, reflecting the island's history.
type Malagasy struct{}

var _ Dialect = (*Malagasy)(nil)

func (m Malagasy) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`malagasy:malagasy;`)
}

func (m Malagasy) Name() string {
	return "malagasy"
}

func (m Malagasy) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"aoka":            Let,
		"maharitra":       Const,
		"raha":            KeywordIf,
		"raha_tsy_izany":  KeywordElse,
		"mandritra":       KeywordWhile,
		"asa":             Function,
		"miverina":        Return,
		"fanombohana":     Main,
		"mamoaka":         Print,
		"rafitra":         Struct,
		"marina":          True,
		"diso":            False,
		"bool":            TypeBool,
		"byte":            TypeByte,
		"float":           TypeFloat,
		"isa_manontolo":   TypeInt,
		"isa_manontolo64": TypeInt64,
		"laharana_litera": TypeString,
		"hadisoana":       TypeError,
	}
}

func (m Malagasy) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
