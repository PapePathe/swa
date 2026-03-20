package lexer

import (
	"regexp"
)

// Country: Ghana
// Ga is a Kwa language spoken by the Ga people in the Greater Accra Region
// of Ghana. it is the principal language of the capital city, Accra, and
// has a significant literary and musical tradition.
type Ga struct{}

var _ Dialect = (*Ga)(nil)

func (m Ga) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`ga:ga;`)
}

func (m Ga) Name() string {
	return "ga"
}

func (m Ga) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ga":       DialectDeclaration,
		"ha":       Let,
		"daa":      Const,
		"kɛ":       KeywordIf,
		"loo":      KeywordElse,
		"mli":      KeywordWhile,
		"nitsumɔ":  Function,
		"ku":       Return,
		"jije":     Main,
		"gba":      Print,
		"nifeemɔ":  Struct,
		"anɔkwale": True,
		"amalɛ":    False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"wiemɔ":    TypeString,
		"tɔmɔ":     TypeError,
		"zero":     Zero,
	}
}

func (m Ga) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
