package lexer

import (
	"regexp"
)

// Country: Uganda
// Nyankole (Runyankore) is a Bantu language spoken by the Banyankore people
// of southwestern Uganda. It is part of the Great Lakes Bantu group and
// is closely related to Rukiga and Runyoro.
type Nyankole struct{}

var _ Dialect = (*Nyankole)(nil)

func (m Nyankole) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`nyankole:nyankole;`)
}

func (m Nyankole) Name() string {
	return "runyankole"
}

func (m Nyankole) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"reka":       Let,
		"gumire":     Const,
		"kuba":       KeywordIf,
		"nari":       KeywordElse,
		"obwo":       KeywordWhile,
		"omulimo":    Function,
		"garura":     Return,
		"tandika":    Main,
		"andika":     Print,
		"mutaratara": Struct,
		"mazima":     True,
		"bishuba":    False,
		"bool":       TypeBool,
		"bauto":      TypeByte,
		"float":      TypeFloat,
		"namba":      TypeInt,
		"namba64":    TypeInt64,
		"ebigambo":   TypeString,
		"ensobi":     TypeError,
	}
}

func (m Nyankole) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
