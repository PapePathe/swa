package lexer

import (
	"regexp"
)

// Country: Ghana
// Kasem is a Gur language spoken in the Upper East Region of Ghana and
// southern Burkina Faso. It is the language of the Kasena people and has
// a unique linguistic structure within the Gur group.
type Kasem struct{}

var _ Dialect = (*Kasem)(nil)

func (m Kasem) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`kasem:kasem;`)
}

func (m Kasem) Name() string {
	return "kasem"
}

func (m Kasem) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"kasem":   DialectDeclaration,
		"de":      Let,
		"dɛidɛi":  Const,
		"tɛ":      KeywordIf,
		"ooo":     KeywordElse,
		"nɛ":      KeywordWhile,
		"tuu":     Function,
		"vuri":    Return,
		"piga":    Main,
		"tora":    Print,
		"liri":    Struct,
		"bobɔ":    True,
		"lanyer":  False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"float":   TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"munyu":   TypeString,
		"tɔmɔ":    TypeError,
		"zero":    Zero,
	}
}

func (m Kasem) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
