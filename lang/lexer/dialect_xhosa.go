package lexer

import (
	"regexp"
)

// Country: South Africa
// isiXhosa is a major Nguni language and one of the eleven official languages
// of South Africa. It is known for its distinctive click sounds and is
// spoken primarily in the Eastern Cape province.
type Xhosa struct{}

var _ Dialect = (*Xhosa)(nil)

func (m Xhosa) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`xhosa:xhosa;`)
}

func (m Xhosa) Name() string {
	return "xhosa"
}

func (m Xhosa) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"makuve":             Let,
		"ngunaphakade":       Const,
		"ukuba":              KeywordIf,
		"okanye":             KeywordElse,
		"ngeli_xesha":        KeywordWhile,
		"msebenzi":           Function,
		"buyisela":           Return,
		"qala":               Main,
		"shicilela":          Print,
		"ulwakhiwo":          Struct,
		"inyani":             True,
		"ubuxoki":            False,
		"bool":               TypeBool,
		"byte":               TypeByte,
		"float":              TypeFloat,
		"inani_elipheleleyo": TypeInt,
		"inani64":            TypeInt64,
		"umbhalo":            TypeString,
		"iphutha":            TypeError,
	}
}

func (m Xhosa) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
