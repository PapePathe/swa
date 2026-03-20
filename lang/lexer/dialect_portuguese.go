package lexer

import (
	"regexp"
)

// Country: Portugal, Brazil, Angola, Mozambique, Cape Verde, etc.
// Portuguese is a Romance language spoken globally. It has a deep-rooted
// presence in Africa, serving as an official language in countries like
// Angola, Mozambique, Guinea-Bissau, and Cape Verde.
type Portuguese struct{}

var _ Dialect = (*Portuguese)(nil)

func (m Portuguese) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialeto:portuguese;`)
}

func (m Portuguese) Name() string {
	return "português"
}

func (m Portuguese) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"seja":       Let,
		"const":      Const,
		"se":         KeywordIf,
		"senao":      KeywordElse,
		"enquanto":   KeywordWhile,
		"funcao":     Function,
		"retornar":   Return,
		"inicio":     Main,
		"imprimir":   Print,
		"estrutura":  Struct,
		"verdadeiro": True,
		"falso":      False,
		"booleano":   TypeBool,
		"byte":       TypeByte,
		"decimal":    TypeFloat,
		"inteiro":    TypeInt,
		"inteiro64":  TypeInt64,
		"texto":      TypeString,
		"erro":       TypeError,
	}
}

func (m Portuguese) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
