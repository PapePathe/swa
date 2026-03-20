package lexer

import (
	"regexp"
)

// Country: Nigeria
// Ibibio is the primary language of the Ibibio people in Akwa Ibom State, South-Southern Nigeria.
// It is a Benue-Congo language that is closely related to Efik and Annang, forming a significant
// cluster in the Niger Delta's linguistic landscape.
type Ibibio struct{}

var _ Dialect = (*Ibibio)(nil)

func (m Ibibio) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`ibibio:ibibio;`)
}

func (m Ibibio) Name() string {
	return "ibibio"
}

func (m Ibibio) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"wek":         Let,
		"ma-pe-lokke": Const,
		"edieke":      KeywordIf,
		"efen":        KeywordElse,
		"dana":        KeywordWhile,
		"kobo":        Function,
		"nyoon":       Return,
		"piilgu":      Main,
		"uminwed":     Print,
		"yub":         Struct,
		"akpaniko":    True,
		"goba":        False,
		"bool":        TypeBool,
		"byte":        TypeByte,
		"decimal":     TypeFloat,
		"namba":       TypeInt,
		"namba64":     TypeInt64,
		"lok-anyen":   TypeString,
		"bal":         TypeError,
	}
}

func (m Ibibio) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
