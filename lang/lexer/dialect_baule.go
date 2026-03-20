package lexer

import (
	"regexp"
)

// Country: Ivory Coast
// Baoulé is a Central Tano language spoken primarily in the central and southern regions
// of Ivory Coast. It is one of the most widely spoken languages in the country and has
// a rich oral tradition and cultural history.
type Baule struct{}

var _ Dialect = (*Baule)(nil)

func (m Baule) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`baule:baule;`)
}

func (m Baule) Name() string {
	return "baoulé"
}

func (m Baule) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"man":      Let,
		"man-kpa":  Const,
		"se":       KeywordIf,
		"se-ka":    KeywordElse,
		"kponin":   KeywordWhile,
		"junman":   Function,
		"sayi":     Return,
		"mo":       Main,
		"kle":      Print,
		"nhyehyee": Struct,
		"nanwle":   True,
		"ngatwa":   False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"desimal":  TypeFloat,
		"lim":      TypeInt,
		"lim64":    TypeInt64,
		"nde":      TypeString,
		"ngasi":    TypeError,
	}
}

func (m Baule) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
