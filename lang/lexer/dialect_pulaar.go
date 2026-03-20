package lexer

import (
	"regexp"
)

// Country: Senegal
// Pulaar is a variety of the Fula language spoken primarily in the Senegal River valley
// region. It is a major Atlantic language with a long history of literacy and is central
// to the identity of the Peul/Haalpulaar'en people.
type Pulaar struct{}

var _ Dialect = (*Pulaar)(nil)

func (m Pulaar) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`pulaar:pulaar;`)
}

func (m Pulaar) Name() string {
	return "pulaar"
}

func (m Pulaar) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"accu":     Let,
		"tabitdum": Const,
		"so":       KeywordIf,
		"so-wonaa": KeywordElse,
		"hono":     KeywordWhile,
		"golle":    Function,
		"rutti":    Return,
		"fuddo":    Main,
		"windu":    Print,
		"mbalndi":  Struct,
		"goonga":   True,
		"fenaande": False,
		"bool":     TypeBool,
		"bayte":    TypeByte,
		"desimal":  TypeFloat,
		"limre":    TypeInt,
		"limre64":  TypeInt64,
		"mbindu":   TypeString,
		"njuumte":  TypeError,
	}
}

func (m Pulaar) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
