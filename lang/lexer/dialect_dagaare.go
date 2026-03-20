package lexer

import (
	"regexp"
)

// Country: Ghana
// Dagaare is a Gur language spoken in the Upper West Region of Ghana and parts
// of Burkina Faso. It is a major language of the Dagaaba people and has several
// distinct dialects across the border areas.
type Dagaare struct{}

var _ Dialect = (*Dagaare)(nil)

func (m Dagaare) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dagaare:dagaare;`)
}

func (m Dagaare) Name() string {
	return "dagaare"
}

func (m Dagaare) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ko":       Let,
		"ma-yɛr":   Const,
		"ka":       KeywordIf,
		"ka-ba":    KeywordElse,
		"kyɛ":      KeywordWhile,
		"toma":     Function,
		"le":       Return,
		"piilo":    Main,
		"sɛge":     Print,
		"nhyehyɛɛ": Struct,
		"sida":     True,
		"ziri":     False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"gom":      TypeString,
		"tuurga":   TypeError,
	}
}

func (m Dagaare) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
