package lexer

import (
	"regexp"
)

// Country: Spain, Latin America, Equatorial Guinea.
// Spanish (Español) is a Romance language that originated in Spain and is
// now the second most spoken native language in the world. It has a
// presence in Africa through Equatorial Guinea and North African enclaves.
type Spanish struct{}

var _ Dialect = (*Spanish)(nil)

func (m Spanish) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialecto:spanish;`)
}

func (m Spanish) Name() string {
	return "español"
}

func (m Spanish) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"declarar":   Let,
		"const":      Const,
		"si":         KeywordIf,
		"sino":       KeywordElse,
		"mientras":   KeywordWhile,
		"fun":        Function,
		"retornar":   Return,
		"inicio":     Main,
		"imprimir":   Print,
		"estructura": Struct,
		"verdadero":  True,
		"falso":      False,
		"booleano":   TypeBool,
		"byte":       TypeByte,
		"decimal":    TypeFloat,
		"entero":     TypeInt,
		"entero64":   TypeInt64,
		"cadena":     TypeString,
		"error":      TypeError,
	}
}

func (m Spanish) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
