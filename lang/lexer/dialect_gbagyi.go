package lexer

import (
	"regexp"
)

// Country: Nigeria
// Gbagyi (Gwari) is a Nupoid language spoken by the Gbagyi people across North-Central
// Nigeria, including the Federal Capital Territory (Abuja). It is one of the most widely
// distributed languages in the Middle Belt.
type Gbagyi struct{}

var _ Dialect = (*Gbagyi)(nil)

func (m Gbagyi) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`gbagyi:gbagyi;`)
}

func (m Gbagyi) Name() string {
	return "gbagyi"
}

func (m Gbagyi) Reserved() map[string]TokenKind {
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

func (m Gbagyi) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
