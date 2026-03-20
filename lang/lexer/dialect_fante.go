package lexer

import (
	"regexp"
)

// Country: Ghana
// Fante is a major dialect of the Akan language spoken primarily in the central and
// western coastal regions of Ghana. It is a key literary and commercial language
// within the Akan cultural cluster.
type Fante struct{}

var _ Dialect = (*Fante)(nil)

func (m Fante) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`fante:fante;`)
}

func (m Fante) Name() string {
	return "fante"
}

func (m Fante) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ma":       Let,
		"pintinn":  Const,
		"sɛ":       KeywordIf,
		"anaasɛ":   KeywordElse,
		"mber":     KeywordWhile,
		"edwuma":   Function,
		"san":      Return,
		"mfitiase": Main,
		"kyerɛw":   Print,
		"nhyehyɛɛ": Struct,
		"nokwar":   True,
		"ator":     False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"nsɛm":     TypeString,
		"mfomso":   TypeError,
	}
}

func (m Fante) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
