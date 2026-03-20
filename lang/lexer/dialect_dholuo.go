package lexer

import (
	"regexp"
)

// Country: Kenya
// Dholuo is a Nilotic language spoken by the Luo people of Kenya and Tanzania.
// It is the most widely spoken Nilotic language in Kenya and is central
// to the cultural and intellectual life of the Lake Victoria region.
type Dholuo struct{}

var _ Dialect = (*Dholuo)(nil)

func (m Dholuo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dholuo:dholuo;`)
}

func (m Dholuo) Name() string {
	return "dholuo"
}

func (m Dholuo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"wek":          Let,
		"ma-ok-lok":    Const,
		"kaponi":       KeywordIf,
		"ka_ok_kamano": KeywordElse,
		"kindego_duto": KeywordWhile,
		"tich":         Function,
		"miyo_dok":     Return,
		"chak":         Main,
		"nyiso":        Print,
		"yuto":         Struct,
		"adier":        True,
		"mieri":        False,
		"bool":         TypeBool,
		"bayiti":       TypeByte,
		"desimal":      TypeFloat,
		"mangima":      TypeInt,
		"mangima64":    TypeInt64,
		"weche":        TypeString,
		"ketho":        TypeError,
	}
}

func (m Dholuo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
