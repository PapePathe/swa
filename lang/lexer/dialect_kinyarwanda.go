package lexer

import (
	"regexp"
)

// Country: Rwanda
// Kinyarwanda is the national language of Rwanda and one of its four official
// languages. It is a Bantu language widely understood throughout the country
// and is key to its social, political, and cultural identity.
type Kinyarwanda struct{}

var _ Dialect = (*Kinyarwanda)(nil)

func (m Kinyarwanda) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`kinyarwanda:kinyarwanda;`)
}

func (m Kinyarwanda) Name() string {
	return "kinyarwanda"
}

func (m Kinyarwanda) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"reka":                Let,
		"udahinduka":          Const,
		"niba":                KeywordIf,
		"ikindi":              KeywordElse,
		"mugihe":              KeywordWhile,
		"porogaramu_ntoya":    Function,
		"tanga":               Return,
		"tangira":             Main,
		"tangaza_amakuru":     Print,
		"imiterere":           Struct,
		"ukuri":               True,
		"ibinyoma":            False,
		"bool":                TypeBool,
		"byte":                TypeByte,
		"float":               TypeFloat,
		"umubare_wuzuye":      TypeInt,
		"umubare_wuzuye64":    TypeInt64,
		"ikurikiranyanyuguti": TypeString,
		"ikosa":               TypeError,
	}
}

func (m Kinyarwanda) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
