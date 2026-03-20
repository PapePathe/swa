package lexer

import (
	"regexp"
)

// Country: Niger
// Zarma is a Songhai language spoken primarily in the south-western regions of Niger.
// It is the second most spoken language in the country and is central to the
// cultural life of the Zarma and Songhai peoples.
type Zarma struct{}

var _ Dialect = (*Zarma)(nil)

func (m Zarma) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`zarma:zarma;`)
}

func (m Zarma) Name() string {
	return "zarma"
}

func (m Zarma) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"nan":     Let,
		"tabat":   Const,
		"da":      KeywordIf,
		"da-ma":   KeywordElse,
		"waati":   KeywordWhile,
		"kaadin":  Function,
		"yeere":   Return,
		"sintin":  Main,
		"kaarye":  Print,
		"laben":   Struct,
		"cimi":    True,
		"tangari": False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"desimal": TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"hantum":  TypeString,
		"daray":   TypeError,
	}
}

func (m Zarma) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
