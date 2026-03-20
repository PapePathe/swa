package lexer

import (
	"regexp"
)

// Country: Ghana
// Gonja is a North Guang language spoken primarily in the Savannah Region
// of Ghana. It is the language of the Gonja people and has a rich oral
// history and cultural heritage dating back to the 16th century.
type Gonja struct{}

var _ Dialect = (*Gonja)(nil)

func (m Gonja) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`gonja:gonja;`)
}

func (m Gonja) Name() string {
	return "gonja"
}

func (m Gonja) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"gonja":       DialectDeclaration,
		"de":          Let,
		"kpa":         Const,
		"mfo":         KeywordIf,
		"ooo":         KeywordElse,
		"be":          KeywordWhile,
		"tumpung":     Function,
		"kpale":       Return,
		"keshishishi": Main,
		"lani":        Print,
		"nkyekyɛle":   Struct,
		"ebase":       True,
		"ebade":       False,
		"bool":        TypeBool,
		"byte":        TypeByte,
		"float":       TypeFloat,
		"lamba":       TypeInt,
		"lamba64":     TypeInt64,
		"muuny":       TypeString,
		"mfonsho":     TypeError,
		"zero":        Zero,
	}
}

func (m Gonja) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
