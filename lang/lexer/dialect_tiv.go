package lexer

import (
	"regexp"
)

// Country: Nigeria
// Tiv is a significant Benue-Congo language spoken by the Tiv people, primarily in Benue and
// Taraba states of North-Central Nigeria. It is known for its unique linguistic structure and
// its central role in the cultural life of the Middle Belt region.
type Tiv struct{}

var _ Dialect = (*Tiv)(nil)

func (m Tiv) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`tiv:tiv;`)
}

func (m Tiv) Name() string {
	return "tiv"
}

func (m Tiv) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"lu":        Let,
		"kpa-vough": Const,
		"aluer":     KeywordIf,
		"kpa":       KeywordElse,
		"shighe-u":  KeywordWhile,
		"tichi":     Function,
		"dwok-cen":  Return,
		"mase":      Main,
		"go-piny":   Print,
		"thing":     Struct,
		"mimi":      True,
		"mimi-ga":   False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"decimal":   TypeFloat,
		"iyenge":    TypeInt,
		"iyenge64":  TypeInt64,
		"mo":        TypeString,
		"bo":        TypeError,
	}
}

func (m Tiv) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
