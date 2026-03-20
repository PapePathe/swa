package lexer

import (
	"regexp"
)

// Country: Eritrea
// Tigrinya is an Ethiopian Semitic language spoken primarily in Eritrea
// and the Tigray Region of Ethiopia. It is a major language of the
// Horn of Africa and uses the centuries-old Ge'ez script.
type Tigrinya struct{}

var _ Dialect = (*Tigrinya)(nil)

func (m Tigrinya) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`tigrinya:tigrinya;`)
}

func (m Tigrinya) Name() string {
	return "tigrinya"
}

func (m Tigrinya) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ይኹን":      Let,
		"ቐዋሚ":      Const,
		"እንተ":      KeywordIf,
		"እንተዘይኮይኑ": KeywordElse,
		"ከሎ":       KeywordWhile,
		"ተግባር":     Function,
		"መልስ":      Return,
		"ጀምር":      Main,
		"ሕተም":      Print,
		"መዋቅር":     Struct,
		"ሓቂ":       True,
		"ሓሶት":      False,
		"bool":     TypeBool,
		"byte":     TypeByte,
		"float":    TypeFloat,
		"ቁጽሪ":      TypeInt,
		"ቁጽሪ64":    TypeInt64,
		"ቃል":       TypeString,
		"ጌጋ":       TypeError,
	}
}

func (m Tigrinya) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
