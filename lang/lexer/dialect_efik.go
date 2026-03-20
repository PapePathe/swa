package lexer

import (
	"regexp"
)

// Country: Nigeria
// Efik is a prestigious Benue-Congo language spoken primarily in Cross River State,
// South-Southern Nigeria. It has historically served as a literary and trade language
// across the Cross River region and is closely related to Ibibio.
type Efik struct{}

var _ Dialect = (*Efik)(nil)

func (m Efik) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`efik:efik;`)
}

func (m Efik) Name() string {
	return "efik"
}

func (m Efik) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"yak":      Let,
		"nsinsi":   Const,
		"edieke":   KeywordIf,
		"efen":     KeywordElse,
		"ke-ini":   KeywordWhile,
		"kobo":     Function,
		"fiak":     Return,
		"mbido":    Main,
		"uminwed":  Print,
		"nhazi":    Struct,
		"akpaniko": True,
		"uwo":      False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"decimal":  TypeFloat,
		"iyenge":   TypeInt,
		"iyenge64": TypeInt64,
		"okwu":     TypeString,
		"bal":      TypeError,
	}
}

func (m Efik) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
