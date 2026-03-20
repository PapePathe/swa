package lexer

import (
	"regexp"
)

// Country: South Africa
// Afrikaans is a West Germanic language that evolved from Dutch 17th-century
// dialects in Southern Africa. It is one of the official languages of
// South Africa and is a major language of culture, science, and media.
type Afrikaans struct{}

var _ Dialect = (*Afrikaans)(nil)

func (m Afrikaans) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`afrikaans:afrikaans;`)
}

func (m Afrikaans) Name() string {
	return "afrikaans"
}

func (m Afrikaans) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"laat":       Let,
		"konstant":   Const,
		"as":         KeywordIf,
		"anders":     KeywordElse,
		"terwyl":     KeywordWhile,
		"funksie":    Function,
		"terugstuur": Return,
		"hoof":       Main,
		"druk":       Print,
		"struktuur":  Struct,
		"waar":       True,
		"onwaar":     False,
		"bool":       TypeBool,
		"greep":      TypeByte,
		"desimaal":   TypeFloat,
		"getal":      TypeInt,
		"getal64":    TypeInt64,
		"teks":       TypeString,
		"fout":       TypeError,
	}
}

func (m Afrikaans) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
