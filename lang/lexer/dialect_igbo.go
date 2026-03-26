package lexer

import (
	"regexp"
)

// Igbo
//
// Country: Nigeria
type Igbo struct{}

var _ Dialect = (*Igbo)(nil)

func (m Igbo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`oluasusu:igbo;`)
}

func (m Igbo) Name() string {
	return "igbo"
}

func (m Igbo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"oluasusu":   DialectDeclaration,
		"ka":         Let,
		"ogidi":      Const,
		"maobu":      KeywordElse,
		"oburu":      KeywordIf,
		"dika":       KeywordWhile,
		"oru":        Function,
		"weghachi":   Return,
		"mbido":      Main,
		"deputa":     Print,
		"nhazi":      Struct,
		"ezigbo":     True,
		"ụgha":       False,
		"ezigha":     TypeBool,
		"mkpụrụ":     TypeByte,
		"ofe":        TypeFloat,
		"onuogugu":   TypeInt,
		"onuogugu64": TypeInt64,
		"ederede":    TypeString,
		"njehie":     TypeError,
		"ogologo":    Variadic,
		"efu":        Zero,
	}
}

func (m Igbo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
