package lexer

import (
	"regexp"
)

// Country: Nigeria
// Izon (Ijaw) is spoken by the Ijaw people, primarily in the Niger Delta region of Nigeria,
// including Bayelsa, Delta, and Rivers states. It is a member of the Ijoid language family
// and is key to the cultural and social life of the coastal communities.
type Izon struct{}

var _ Dialect = (*Izon)(nil)

func (m Izon) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`izon:izon;`)
}

func (m Izon) Name() string {
	return "izon"
}

func (m Izon) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"bi":      Let,
		"tabat":   Const,
		"be":      KeywordIf,
		"be-mbi":  KeywordElse,
		"egbe-ni": KeywordWhile,
		"baara":   Function,
		"kan":     Return,
		"start":   Main,
		"binda":   Print,
		"egbe":    Struct,
		"ijaw":    True,
		"ijaw-ga": False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"decimal": TypeFloat,
		"namba":   TypeInt,
		"namba64": TypeInt64,
		"mo":      TypeString,
		"bal":     TypeError,
	}
}

func (m Izon) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
