package lexer

import (
	"regexp"
)

// Country: Nigeria
// Edo, also known as Bini, is the language of the Edo people of Edo State, South-Western Nigeria.
// It was the primary language of the ancient Benin Empire and remains a prestigious language
// of history, art, and tradition in the region.
type Edo struct{}

var _ Dialect = (*Edo)(nil)

func (m Edo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`edo:edo;`)
}

func (m Edo) Name() string {
	return "edo"
}

func (m Edo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"gi":        Let,
		"no-khun":   Const,
		"ghe-okhe":  KeywordIf,
		"ofa":       KeywordElse,
		"eghe-ne":   KeywordWhile,
		"wia":       Function,
		"rhivwin":   Return,
		"fara":      Main,
		"gben":      Print,
		"unhunvwun": Struct,
		"emwanta":   True,
		"emwanta-i": False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"decimal":   TypeFloat,
		"namba":     TypeInt,
		"namba64":   TypeInt64,
		"eta":       TypeString,
		"obo":       TypeError,
	}
}

func (m Edo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
