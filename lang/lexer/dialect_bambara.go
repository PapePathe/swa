package lexer

import (
	"regexp"
)

// Country: Mali
// Bambara (Bamanankan) is the primary lingua franca of Mali, spoken by the vast
// majority of the population. It is a Mande language with a rich tradition of
// storytelling and serves as a major language of trade and administration.
type Bambara struct{}

var _ Dialect = (*Bambara)(nil)

func (m Bambara) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`bm:bm;`)
}

func (m Bambara) Name() string {
	return "bamanankan"
}

func (m Bambara) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"a_to":      Let,
		"jigi":      Const,
		"ni":        KeywordIf,
		"o_tɛ_o_ye": KeywordElse,
		"tuma_min":  KeywordWhile,
		"baara":     Function,
		"kà_segi":   Return,
		"labɛn":     Main,
		"sɛbɛn":     Print,
		"ton":       Struct,
		"tiɲɛ":      True,
		"nkalon":    False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"decimal":   TypeFloat,
		"limo":      TypeInt,
		"limo64":    TypeInt64,
		"mbinde":    TypeString,
		"filila":    TypeError,
	}
}

func (m Bambara) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
