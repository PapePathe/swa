package lexer

import (
	"regexp"
)

// Country: Ghana
// Ewe is a Gbe language spoken in southeastern Ghana, southern Togo, and Benin.
// It has a significant literary tradition and is a major language of
// education and broadcast media in the region.
type Ewe struct{}

var _ Dialect = (*Ewe)(nil)

func (m Ewe) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`ewe:ewe;`)
}

func (m Ewe) Name() string {
	return "ewe"
}

func (m Ewe) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"na":          Let,
		"ma-trɔ-trɔ":  Const,
		"ne":          KeywordIf,
		"ne-menye":    KeywordElse,
		"esi":         KeywordWhile,
		"dɔ":          Function,
		"trɔ":         Return,
		"gɔme":        Main,
		"ŋlɔ":         Print,
		"ɖoɖo":        Struct,
		"nyateƒe":     True,
		"alakpa":      False,
		"bool":        TypeBool,
		"byte":        TypeByte,
		"desimal":     TypeFloat,
		"akɔntaxlẽ":   TypeInt,
		"akɔntaxlẽ64": TypeInt64,
		"nya":         TypeString,
		"vodada":      TypeError,
	}
}

func (m Ewe) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
