package lexer

import (
	"regexp"
)

// Country: Democratic Republic of the Congo
// Lingala is a Bantu language spoken as a lingua franca throughout a large part
// of the Democratic Republic of the Congo and the Republic of the Congo. It is the
// language of modern music (Soukous) and massive urban centers like Kinshasa.
type Lingala struct{}

var _ Dialect = (*Lingala)(nil)

func (m Lingala) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`lingala:lingala;`)
}

func (m Lingala) Name() string {
	return "lingala"
}

func (m Lingala) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"tika":      Let,
		"ngwi":      Const,
		"soki":      KeywordIf,
		"te-soki":   KeywordElse,
		"ntango":    KeywordWhile,
		"mosala":    Function,
		"zongisa":   Return,
		"banda":     Main,
		"lakisa":    Print,
		"molongo":   Struct,
		"solo":      True,
		"lokuta":    False,
		"bool":      TypeBool,
		"byte":      TypeByte,
		"desimal":   TypeFloat,
		"motango":   TypeInt,
		"motango64": TypeInt64,
		"maloba":    TypeString,
		"libunga":   TypeError,
	}
}

func (m Lingala) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
