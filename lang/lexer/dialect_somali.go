package lexer

import (
	"regexp"
)

// Country: Somalia
// Somali is a Cushitic language spoken as a national language in Somalia,
// as well as in Djibouti, Ethiopia, and Kenya. It has a rich poetic
// tradition and is the most documented Cushitic language.
type Somali struct{}

var _ Dialect = (*Somali)(nil)

func (m Somali) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`somali:somali;`)
}

func (m Somali) Name() string {
	return "somali"
}

func (m Somali) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ha":          Let,
		"joogto":      Const,
		"haddii":      KeywordIf,
		"haddii_kale": KeywordElse,
		"halka":       KeywordWhile,
		"shaqo":       Function,
		"celi":        Return,
		"bilow":       Main,
		"daabac":      Print,
		"dhisme":      Struct,
		"run":         True,
		"been":        False,
		"bool":        TypeBool,
		"bayt":        TypeByte,
		"tobanle":     TypeFloat,
		"tiro":        TypeInt,
		"tiro64":      TypeInt64,
		"qoraal":      TypeString,
		"qalad":       TypeError,
	}
}

func (m Somali) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
