package lexer

import (
	"regexp"
)

// Country: Senegal
// Diola (Jola) is a cluster of languages spoken in the Casamance region of Senegal,
// Gambia, and Guinea-Bissau. It belongs to the Bak branch of the Niger-Congo family
// and is central to the agrarian culture of the region.
type Diola struct{}

var _ Dialect = (*Diola)(nil)

func (m Diola) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`diola:diola;`)
}

func (m Diola) Name() string {
	return "diola"
}

func (m Diola) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"kanyen":  Let,
		"xeet":    Const,
		"be":      KeywordIf,
		"be-mbi":  KeywordElse,
		"jamook":  KeywordWhile,
		"baara":   Function,
		"kan":     Return,
		"tambali": Main,
		"binda":   Print,
		"kulan":   Struct,
		"kaan":    True,
		"fen":     False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"desimal": TypeFloat,
		"lim":     TypeInt,
		"lim64":   TypeInt64,
		"mbinda":  TypeString,
		"njuumte": TypeError,
	}
}

func (m Diola) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
