package lexer

import (
	"regexp"
)

// Country: Nigeria
// Idoma is a Benue-Congo language spoken by the Idoma people in Benue State, North-Central Nigeria.
// It has several dialects and is a major language of the Middle Belt, reflecting a rich cultural heritage.
type Idoma struct{}

var _ Dialect = (*Idoma)(nil)

func (m Idoma) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`idoma:idoma;`)
}

func (m Idoma) Name() string {
	return "idoma"
}

func (m Idoma) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"lu":        Let,
		"tabat":     Const,
		"edieke":    KeywordIf,
		"efen":      KeywordElse,
		"dana":      KeywordWhile,
		"kobo":      Function,
		"nyoon":     Return,
		"piilgu":    Main,
		"uminwed":   Print,
		"yub":       Struct,
		"okwei":     True,
		"ekemgbe":   False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"decimal":   TypeFloat,
		"namba":     TypeInt,
		"namba64":   TypeInt64,
		"lok-anyen": TypeString,
		"bal":       TypeError,
	}
}

func (m Idoma) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
