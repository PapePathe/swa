package lexer

import (
	"regexp"
)

// Country: Nigeria
// Igbo is a major Nigerian language spoken primarily in the South-Eastern region by the
// Igbo people. It is a tonal language in the Benue–Congo branch of the Niger-Congo family
// and is central to the cultural and commercial identity of Igboland.
type Igbo struct{}

var _ Dialect = (*Igbo)(nil)

func (m Igbo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`igbo:igbo;`)
}

func (m Igbo) Name() string {
	return "igbo"
}

func (m Igbo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ka":       Let,
		"kwụm":     Const,
		"ọbụrụ":    KeywordIf,
		"ọbụrụ-na": KeywordElse,
		"mgbe":     KeywordWhile,
		"ọrụ":      Function,
		"lọghachi": Return,
		"mbido":    Main,
		"pụta":     Print,
		"nhazi":    Struct,
		"ezigbo":   True,
		"ụgha":     False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"desimal":  TypeFloat,
		"nọmba":    TypeInt,
		"nọmba64":  TypeInt64,
		"okwu":     TypeString,
		"njehie":   TypeError,
	}
}

func (m Igbo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
