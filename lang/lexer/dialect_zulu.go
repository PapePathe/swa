package lexer

import (
	"regexp"
)

// Country: South Africa
// isiZulu is a major Nguni language and the most widely spoken home language
// in South Africa. It is an official language with a rich literary
// tradition and is central to the history of the Zulu Kingdom.
type Zulu struct{}

var _ Dialect = (*Zulu)(nil)

func (m Zulu) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`zulu:zulu;`)
}

func (m Zulu) Name() string {
	return "zulu"
}

func (m Zulu) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"vumela":     Let,
		"nqumile":    Const,
		"uma":        KeywordIf,
		"kungenjalo": KeywordElse,
		"ngenkathi":  KeywordWhile,
		"umsebenzi":  Function,
		"buyisa":     Return,
		"qala":       Main,
		"bhala":      Print,
		"isakhiwo":   Struct,
		"iqiniso":    True,
		"amanga":     False,
		"bool":       TypeBool,
		"byte":       TypeByte,
		"decimal":    TypeFloat,
		"inombolo":   TypeInt,
		"inombolo64": TypeInt64,
		"umbhalo":    TypeString,
		"iphutha":    TypeError,
	}
}

func (m Zulu) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
