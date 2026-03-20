package lexer

import (
	"regexp"
)

// Country: Nigeria
// Yoruba is one of the three major languages of Nigeria, primarily spoken in the South-Western
// region. It is a tonal language belonging to the Volta-Niger branch of the Niger-Congo family,
// with a rich tradition of oral literature and poetry.
type Yoruba struct{}

var _ Dialect = (*Yoruba)(nil)

func (m Yoruba) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`yoruba:yoruba;`)
}

func (m Yoruba) Name() string {
	return "yoruba"
}

func (m Yoruba) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"jẹ́":      Let,
		"dandan":   Const,
		"bí":       KeywordIf,
		"míràn":    KeywordElse,
		"nígbàtí":  KeywordWhile,
		"iṣẹ́":     Function,
		"padà":     Return,
		"bẹ̀rẹ̀":   Main,
		"tẹ̀jáde":  Print,
		"ètò":      Struct,
		"òtítọ́":   True,
		"irọ́":     False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"nọ́mbà":   TypeInt,
		"nọ́mbà64": TypeInt64,
		"ọ̀rọ̀":    TypeString,
		"àṣìṣe":    TypeError,
	}
}

func (m Yoruba) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
