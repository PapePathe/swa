package lexer

import (
	"regexp"
)

// Country: Sierra Leone
// Krio is the lingua franca of Sierra Leone, spoken by the vast majority of the
// population either as a first or second language. It is an English-based
// creole with significant influences from African languages like Yoruba and Mende.
type Krio struct{}

var _ Dialect = (*Krio)(nil)

func (m Krio) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`krio:krio;`)
}

func (m Krio) Name() string {
	return "krio"
}

func (m Krio) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"leh":     Let,
		"tranga":  Const,
		"if":      KeywordIf,
		"o-els":   KeywordElse,
		"we":      KeywordWhile,
		"wok":     Function,
		"go-bak":  Return,
		"fos":     Main,
		"rite":    Print,
		"betteh":  Struct,
		"tru":     True,
		"lye":     False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"float":   TypeFloat,
		"namba":   TypeInt,
		"namba64": TypeInt64,
		"tok":     TypeString,
		"mistake": TypeError,
	}
}

func (m Krio) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
