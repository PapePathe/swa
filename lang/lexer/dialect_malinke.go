package lexer

import (
	"fmt"
	"regexp"
)

type Malinke struct{}

var _ Dialect = (*Malinke)(nil)

func (Malinke) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile("kan:malinke;")
}

func (Malinke) Name() string {
	return "malinke"
}

func (m Malinke) Error(key string, args ...any) error {
	return fmt.Errorf("Not implemented")
}

func (m Malinke) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"tintin":    Const,
		"kan":       DialectDeclaration,
		"baara":     Function,
		"wala":      KeywordElse,
		"ni":        KeywordIf,
		"tuma":      KeywordWhile,
		"atö":       Let,
		"daminen":   Main,
		"yira":      Print,
		"segin":     Return,
		"joyoro":    Struct,
		"jatelen":   TypeFloat,
		"jate":      TypeInt,
		"jate64":    TypeInt64,
		"seben":     TypeString,
		"fili":      TypeError,
		"caaman":    Variadic,
		"foy":       Zero,
		"tinye":     True,
		"wouya":     False,
		"tinyejate": TypeBool,
	}
}
