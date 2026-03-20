package lexer

import (
	"regexp"
)

// Country: Lesotho
// Sesotho (Southern Sotho) is the national language of Lesotho and one
// of the eleven official languages of South Africa. It is a Sotho-Tswana
// language with a rich literary history and vast cultural significance.
type Sesotho struct{}

var _ Dialect = (*Sesotho)(nil)

func (m Sesotho) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`sesotho:sesotho;`)
}

func (m Sesotho) Name() string {
	return "sesotho"
}

func (m Sesotho) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"a_re":        Let,
		"sa_fetoleng": Const,
		"haeba":       KeywordIf,
		"kapa":        KeywordElse,
		"ha":          KeywordWhile,
		"mosebetsi":   Function,
		"khutlisa":    Return,
		"qalo":        Main,
		"ngola":       Print,
		"sebopeho":    Struct,
		"nete":        True,
		"makaoma":     False,
		"bool":        TypeBool,
		"byte":        TypeByte,
		"float":       TypeFloat,
		"nomoro":      TypeInt,
		"nomoro64":    TypeInt64,
		"mantsoe":     TypeString,
		"phoso":       TypeError,
	}
}

func (m Sesotho) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
