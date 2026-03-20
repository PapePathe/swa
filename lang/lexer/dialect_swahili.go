package lexer

import (
	"regexp"
)

// Country: Tanzania & Kenya
// Swahili (Kiswahili) is the most widely spoken language in sub-Saharan Africa.
// It is a Bantu language with significant Arabic influence and serves as a
// major lingua franca across East and Central Africa, from the coast to the interior.
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
		"wacha":    Let,
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
		"uwongo":   False,
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
	// Fallback to English errors for now as they are technical compiler messages
	return English{}.Error(key, args...)
}
