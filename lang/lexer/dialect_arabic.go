package lexer

import (
	"regexp"
)

// Country: North Africa (Egypt, Algeria, Morocco, etc.)
// Arabic is a major Afroasiatic language spoken across North Africa and
// the Middle East. It is a global language of science, religion, and literature
// with numerous regional dialects and a standardized Modern Standard form.
type Arabic struct{}

var _ Dialect = (*Arabic)(nil)

func (m Arabic) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`arabic:arabic;`)
}

func (m Arabic) Name() string {
	return "arabic"
}

func (m Arabic) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"دع":       Let,
		"ثابت":     Const,
		"إذا":      KeywordIf,
		"وإلا":     KeywordElse,
		"بينما":    KeywordWhile,
		"دالة":     Function,
		"أعد":      Return,
		"رئيسي":    Main,
		"اطبع":     Print,
		"هيكل":     Struct,
		"صحيح":     True,
		"خطأ":      False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"عدد":      TypeInt,
		"عدد64":    TypeInt64,
		"نص":       TypeString,
		"خطأ_نوعي": TypeError,
	}
}

func (m Arabic) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
