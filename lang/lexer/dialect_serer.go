package lexer

import (
	"regexp"
)

// Country: Senegal
// Serer is a West Atlantic language spoken by the Serer people of Senegal and Gambia.
// It is notable for its complex noun-class system and its close relationship with
// the Wolof and Fula languages.
type Serer struct{}

var _ Dialect = (*Serer)(nil)

func (m Serer) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`serer:serer;`)
}

func (m Serer) Name() string {
	return "seereer"
}

func (m Serer) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"nuun":      Let,
		"fena":      Const,
		"le":        KeywordIf,
		"le-mbaa":   KeywordElse,
		"ga-mook":   KeywordWhile,
		"baara":     Function,
		"fodik":     Return,
		"debbi":     Main,
		"binda":     Print,
		"caakaar":   Struct,
		"kaan":      True,
		"o-ref-waa": False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"desimal":   TypeFloat,
		"limo":      TypeInt,
		"limo64":    TypeInt64,
		"mbinda":    TypeString,
		"njuumte":   TypeError,
	}
}

func (m Serer) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
