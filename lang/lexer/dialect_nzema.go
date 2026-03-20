package lexer

import (
	"regexp"
)

// Country: Ghana
// Nzema is a Central Tano language spoken primarily in the Western Region
// of Ghana and eastern Ivory Coast. It is the language of the Nzema people
// and has a significant presence in the coastal areas.
type Nzema struct{}

var _ Dialect = (*Nzema)(nil)

func (m Nzema) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`nzema:nzema;`)
}

func (m Nzema) Name() string {
	return "nzema"
}

func (m Nzema) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"nzema":      DialectDeclaration,
		"ma":         Let,
		"tendenle":   Const,
		"saa":        KeywordIf,
		"mɔɔ":        KeywordElse,
		"mekɛ":       KeywordWhile,
		"gyima":      Function,
		"sia":        Return,
		"mɔlebɛbo":   Main,
		"kile":       Print,
		"ngyehyɛleɛ": Struct,
		"nɔhalɛ":     True,
		"adaada":     False,
		"bool":       TypeBool,
		"byte":       TypeByte,
		"float":      TypeFloat,
		"namba":      TypeInt,
		"namba64":    TypeInt64,
		"edwɛkɛ":     TypeString,
		"mfomsoɔ":    TypeError,
		"zero":       Zero,
	}
}

func (m Nzema) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
