package lexer

import (
	"regexp"
)

// Country: Ghana
// Dagbani is a Gur language spoken primarily by the Dagomba people in Northern
// Ghana. It is a major language of administration and commerce in the region
// and is closely related to Mampruli and Nanuni.
type Dagbani struct{}

var _ Dialect = (*Dagbani)(nil)

func (m Dagbani) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dagbani:dagbani;`)
}

func (m Dagbani) Name() string {
	return "dagbani"
}

func (m Dagbani) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"che":       Let,
		"ma-yu-yɛr": Const,
		"di":        KeywordIf,
		"di-ba":     KeywordElse,
		"ka":        KeywordWhile,
		"tuma":      Function,
		"labi":      Return,
		"piilgu":    Main,
		"sabbi":     Print,
		"nhyehyɛɛ":  Struct,
		"yɛlimaŋli": True,
		"ʒiri":      False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"float":     TypeFloat,
		"namba":     TypeInt,
		"namba64":   TypeInt64,
		"yɛltɔɣa":   TypeString,
		"taali":     TypeError,
	}
}

func (m Dagbani) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
