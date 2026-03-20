package lexer

import (
	"regexp"
)

// Country: Benin
// Fon is a Gbe language spoken primarily in Benin and parts of Togo and Nigeria. It was
// the primary language of the Kingdom of Dahomey and remains a major cultural and
// communicative force in modern Benin.
type Fon struct{}

var _ Dialect = (*Fon)(nil)

func (m Fon) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`fon:fon;`)
}

func (m Fon) Name() string {
	return "fon"
}

func (m Fon) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"nyo":        Let,
		"akwe":       Const,
		"nu":         KeywordIf,
		"nu-ma":      KeywordElse,
		"hwenu":      KeywordWhile,
		"azo":        Function,
		"leko":       Return,
		"be":         Main,
		"xle":        Print,
		"nhyehyee":   Struct,
		"nugbo":      True,
		"vanyavanya": False,
		"bool":       TypeBool,
		"byte":       TypeByte,
		"desimal":    TypeFloat,
		"limo":       TypeInt,
		"limo64":     TypeInt64,
		"wen":        TypeString,
		"afo":        TypeError,
	}
}

func (m Fon) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
