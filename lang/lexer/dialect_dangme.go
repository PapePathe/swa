package lexer

import (
	"regexp"
)

// Country: Ghana
// Dangme (Adangme) is a Kwa language spoken in south-eastern Ghana, by the
// Dangme people. It is closely related to Ga and has several dialects
// including Ada, Krobo, and Ningo.
type Dangme struct{}

var _ Dialect = (*Dangme)(nil)

func (m Dangme) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dangme:dangme;`)
}

func (m Dangme) Name() string {
	return "dangme"
}

func (m Dangme) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"dangme":   DialectDeclaration,
		"ha":       Let,
		"kpee":     Const,
		"ke":       KeywordIf,
		"loo":      KeywordElse,
		"be":       KeywordWhile,
		"nitsumi":  Function,
		"kpale":    Return,
		"sisije":   Main,
		"de":       Print,
		"nifeemi":  Struct,
		"anɔkwale": True,
		"lakpa":    False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"namba":    TypeInt,
		"namba64":  TypeInt64,
		"munyu":    TypeString,
		"tɔmi":     TypeError,
		"zero":     Zero,
	}
}

func (m Dangme) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
