package lexer

import (
	"regexp"
)

// Country: Ghana
// Twi (Asante and Akuapem) is the most widely spoken dialect of the Akan language
// in Ghana. It serves as a major lingua franca across the country and has a
// profound influence on Ghanaian music, media, and daily life.
type Twi struct{}

var _ Dialect = (*Twi)(nil)

func (m Twi) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`twi:twi;`)
}

func (m Twi) Name() string {
	return "twi"
}

func (m Twi) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ma":        Let,
		"tumi":      Const,
		"sɛ":        KeywordIf,
		"anyɛ_sa_a": KeywordElse,
		"berɛ_a":    KeywordWhile,
		"dwumadie":  Function,
		"san":       Return,
		"mfitiase":  Main,
		"tintim":    Print,
		"nhyehyɛe":  Struct,
		"nokorɛ":    True,
		"atorɔ":     False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"float":     TypeFloat,
		"intiga":    TypeInt,
		"intiga64":  TypeInt64,
		"ahoma":     TypeString,
		"mfomsoɔ":   TypeError,
	}
}

func (m Twi) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
