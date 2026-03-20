package lexer

import (
	"regexp"
)

// Country: Nigeria
// Nupe is a Volta-Niger language spoken primarily by the Nupe people in Niger, Kwara, and
// Kogi states of North-Central Nigeria. It has major dialects like Nupe-Central and is
// historically significant in the Niger-Benue confluence area.
type Nupe struct{}

var _ Dialect = (*Nupe)(nil)

func (m Nupe) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`nupe:nupe;`)
}

func (m Nupe) Name() string {
	return "nupe"
}

func (m Nupe) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"kba":     Let,
		"tabat":   Const,
		"ni":      KeywordIf,
		"gade":    KeywordElse,
		"waati":   KeywordWhile,
		"tuma":    Function,
		"labi":    Return,
		"start":   Main,
		"sabbi":   Print,
		"rny":     Struct,
		"cimi":    True,
		"cimi-ga": False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"decimal": TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"mo":      TypeString,
		"filila":  TypeError,
	}
}

func (m Nupe) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
