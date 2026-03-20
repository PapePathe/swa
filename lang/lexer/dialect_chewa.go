package lexer

import (
	"regexp"
)

// Country: Malawi
// Chewa (Nyanja) is a Bantu language and the national language of Malawi.
// It is also widely spoken in Zambia, Mozambique, and Zimbabwe, serving
// as a major language of communication and culture in South-Central Africa.
type Chewa struct{}

var _ Dialect = (*Chewa)(nil)

func (m Chewa) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`chewa:chewa;`)
}

func (m Chewa) Name() string {
	return "chewa"
}

func (m Chewa) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"lekani":     Let,
		"osasintha":  Const,
		"ngati":      KeywordIf,
		"kapena":     KeywordElse,
		"pamene":     KeywordWhile,
		"ntchito":    Function,
		"bweretsani": Return,
		"pachiyambi": Main,
		"lemba":      Print,
		"mapangidwe": Struct,
		"zoona":      True,
		"bodza":      False,
		"bool":       TypeBool,
		"bayiti":     TypeByte,
		"float":      TypeFloat,
		"nambala":    TypeInt,
		"nambala64":  TypeInt64,
		"mawu":       TypeString,
		"cholakwika": TypeError,
	}
}

func (m Chewa) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
