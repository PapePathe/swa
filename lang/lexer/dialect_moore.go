package lexer

import (
	"regexp"
)

// Country: Burkina Faso
// Mòoré is a Gur language spoken primarily in Burkina Faso by the Mossi people.
// It is the most widely spoken language in the country and is central to its
// national identity and cultural heritage.
type Moore struct{}

var _ Dialect = (*Moore)(nil)

func (m Moore) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`moore:moore;`)
}

func (m Moore) Name() string {
	return "moore"
}

func (m Moore) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"bas":      Let,
		"sida":     Const,
		"sane":     KeywordIf,
		"tɩ-ya":    KeywordElse,
		"wakat":    KeywordWhile,
		"tʋʋma":    Function,
		"lebg":     Return,
		"sɩngre":   Main,
		"gʋls":     Print,
		"teelgo":   Struct,
		"sɩd-sɩda": True,
		"ziri":     False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"gomde":    TypeString,
		"tuurga":   TypeError,
	}
}

func (m Moore) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
