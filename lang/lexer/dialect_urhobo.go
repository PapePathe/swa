package lexer

import (
	"regexp"
)

// Country: Nigeria
// Urhobo is an Edoid language spoken by the Urhobo people of Delta State in South-Southern Nigeria.
// It is the major language of the Western Niger Delta and is central to the identity and
// cultural practices of its speakers.
type Urhobo struct{}

var _ Dialect = (*Urhobo)(nil)

func (m Urhobo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`urhobo:urhobo;`)
}

func (m Urhobo) Name() string {
	return "urhobo"
}

func (m Urhobo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"gi":      Let,
		"tabat":   Const,
		"si-be":   KeywordIf,
		"emu-ofa": KeywordElse,
		"while":   KeywordWhile,
		"wia":     Function,
		"rhivwin": Return,
		"start":   Main,
		"si":      Print,
		"thing":   Struct,
		"uyota":   True,
		"efian":   False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"decimal": TypeFloat,
		"namba":   TypeInt,
		"namba64": TypeInt64,
		"eta":     TypeString,
		"obo":     TypeError,
	}
}

func (m Urhobo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
