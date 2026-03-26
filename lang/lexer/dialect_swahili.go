package lexer

import (
	"regexp"
)

// Country: Tanzania & Kenya
type Swahili struct{}

var _ Dialect = (*Swahili)(nil)

func (m Swahili) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`swahili:swahili;`)
}

func (m Swahili) Name() string {
	return "swahili"
}

func (m Swahili) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"acha":     Let,
		"thabiti":  Const,
		"kama":     KeywordIf,
		"sivyo":    KeywordElse,
		"wakati":   KeywordWhile,
		"kazi":     Function,
		"rudisha":  Return,
		"anza":     Main,
		"chapisha": Print,
		"muundo":   Struct,
		"kweli":    True,
		"uongo":    False,
		"bool":     TypeBool,
		"baiti":    TypeByte,
		"desimali": TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"maneno":   TypeString,
		"kosa":     TypeError,
	}
}

func (m Swahili) Error(key string, args ...any) error {
	// FIXME translate error messages to swahili
	return English{}.Error(key, args...)
}
