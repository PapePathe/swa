package lexer

import (
	"regexp"
)

// Country: Nigeria
// Hausa is the most widely spoken language in Northern Nigeria and across much of West Africa.
// It is a Chadic language within the Afroasiatic family, serving as a primary language of
// trade, commerce, and Islamic education in the Sahel region.
type Hausa struct{}

var _ Dialect = (*Hausa)(nil)

func (m Hausa) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`hausa:hausa;`)
}

func (m Hausa) Name() string {
	return "hausa"
}

func (m Hausa) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"bari":       Let,
		"tabbatacce": Const,
		"idan":       KeywordIf,
		"inbahakaba": KeywordElse,
		"yayinda":    KeywordWhile,
		"aiki":       Function,
		"koma":       Return,
		"fara":       Main,
		"buga":       Print,
		"tsarin":     Struct,
		"gaskiya":    True,
		"karya":      False,
		"bool":       TypeBool,
		"baiti":      TypeByte,
		"decimal":    TypeFloat,
		"lamba":      TypeInt,
		"lamba64":    TypeInt64,
		"rubutu":     TypeString,
		"kuskure":    TypeError,
	}
}

func (m Hausa) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
