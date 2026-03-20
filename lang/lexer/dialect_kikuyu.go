package lexer

import (
	"regexp"
)

// Country: Kenya
// Kikuyu (Gikuyu) is a major Bantu language spoken primarily by the Kikuyu
// people, the largest ethnic group in Kenya. It is a key language of
// central Kenya and has a rich tradition of oral literature and resistance.
type Kikuyu struct{}

var _ Dialect = (*Kikuyu)(nil)

func (m Kikuyu) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`kikuyu:kikuyu;`)
}

func (m Kikuyu) Name() string {
	return "gikuyu"
}

func (m Kikuyu) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"reke":         Let,
		"ma-ti-chenji": Const,
		"angikoruo":    KeywordIf,
		"kana":         KeywordElse,
		"riri":         KeywordWhile,
		"wira":         Function,
		"cokania":      Return,
		"ambiriria":    Main,
		"andika":       Print,
		"mutaratara":   Struct,
		"ma":           True,
		"maheni":       False,
		"bool":         TypeBool,
		"baite":        TypeByte,
		"float":        TypeFloat,
		"namba":        TypeInt,
		"namba64":      TypeInt64,
		"kiugo":        TypeString,
		"ihitia":       TypeError,
	}
}

func (m Kikuyu) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
