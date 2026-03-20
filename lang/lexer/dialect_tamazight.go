package lexer

import (
	"regexp"
)

// Country: North Africa (Morocco, Algeria)
// Tamazight (Berber) is a cluster of Afroasiatic languages native to North
// Africa. It has been spoken in the region for millennia and is an
// official language in Morocco and Algeria, with its own unique script (Tifinagh).
type Tamazight struct{}

var _ Dialect = (*Tamazight)(nil)

func (m Tamazight) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`tamazight:tamazight;`)
}

func (m Tamazight) Name() string {
	return "tamazight"
}

func (m Tamazight) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ay":      Let,
		"itabten": Const,
		"iy":      KeywordIf,
		"neg":     KeywordElse,
		"skud":    KeywordWhile,
		"tasɣent": Function,
		"rar":     Return,
		"agejdan": Main,
		"sigg":    Print,
		"talɣa":   Struct,
		"tidett":  True,
		"arkis":   False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"float":   TypeFloat,
		"amḍan":   TypeInt,
		"amḍan64": TypeInt64,
		"aḍris":   TypeString,
		"tuccḍa":  TypeError,
	}
}

func (m Tamazight) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
