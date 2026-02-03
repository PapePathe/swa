package lexer

import (
	"fmt"
	"regexp"
)

type Malinke struct{}

var _ Dialect = (*Malinke)(nil)

func (Malinke) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:malinke;`)
}

func (m Malinke) Error(key string, args ...any) error {
	return fmt.Errorf("Not implemented")
}

func (m Malinke) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ni":     KeywordIf,
		"nii":    KeywordElse,
		"struct": Struct,
		"let":    Let,
		"const":  Const,
		"fèndo":  TypeInt,
		"erreur": TypeError,
	}
}
