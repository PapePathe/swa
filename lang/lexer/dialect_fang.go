package lexer

import (
	"regexp"
)

// Country: Gabon
// Fang is a major Bantu language spoken in Gabon, Equatorial Guinea, and southern
// Cameroon. It is the language of the Fang people and is central to the
// oral history and cultural identity of the region.
type Fang struct{}

var _ Dialect = (*Fang)(nil)

func (m Fang) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`fang:fang;`)
}

func (m Fang) Name() string {
	return "fang"
}

func (m Fang) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"vae":     Let,
		"ngo":     Const,
		"nge":     KeywordIf,
		"nge-ki":  KeywordElse,
		"eyong":   KeywordWhile,
		"ese":     Function,
		"bulan":   Return,
		"tat":     Main,
		"tsini":   Print,
		"caakaar": Struct,
		"afiri":   True,
		"vul":     False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"desimal": TypeFloat,
		"namba":   TypeInt,
		"namba64": TypeInt64,
		"nkobo":   TypeString,
		"bitsut":  TypeError,
	}
}

func (m Fang) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
