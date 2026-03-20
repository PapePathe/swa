package lexer

import (
	"regexp"
)

// Country: Central African Republic
// Sango is the primary lingua franca and a national language of the Central
// African Republic. It is a Ngbandi-based creole that facilitates communication
// across the country's diverse ethnic and linguistic groups.
type Sango struct{}

var _ Dialect = (*Sango)(nil)

func (m Sango) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`sango:sango;`)
}

func (m Sango) Name() string {
	return "sango"
}

func (m Sango) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"zia":           Let,
		"nduru":         Const,
		"tongana":       KeywordIf,
		"wala":          KeywordElse,
		"na-ngo":        KeywordWhile,
		"kua":           Function,
		"kiri":          Return,
		"to-ndo":        Main,
		"fa-ra":         Print,
		"da":            Struct,
		"tene-ti-nduru": True,
		"wataka":        False,
		"bool":          TypeBool,
		"byte":          TypeByte,
		"float":         TypeFloat,
		"wongo":         TypeInt,
		"wongo64":       TypeInt64,
		"tene":          TypeString,
		"futi":          TypeError,
	}
}

func (m Sango) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
