package lexer

import (
	"regexp"
)

// Country: Ethiopia
// Amharic is an Ethiopian Semitic language and the primary working language
// of the Ethiopian government. It has its own unique script (Ge'ez) and
// is the second most spoken Semitic language in the world after Arabic.
type Amharic struct{}

var _ Dialect = (*Amharic)(nil)

func (m Amharic) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`amharic:amharic;`)
}

func (m Amharic) Name() string {
	return "amharic"
}

func (m Amharic) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ይሁን":   Let,
		"ቋሚ":    Const,
		"ከሆነ":   KeywordIf,
		"ካልሆነ":  KeywordElse,
		"ሳለ":    KeywordWhile,
		"ተግባር":  Function,
		"መልስ":   Return,
		"ጀምር":   Main,
		"አትም":   Print,
		"መዋቅር":  Struct,
		"እውነት":  True,
		"ውሸት":   False,
		"bool":  TypeBool,
		"byte":  TypeByte,
		"float": TypeFloat,
		"ቁጥር":   TypeInt,
		"ቁጥር64": TypeInt64,
		"ቃል":    TypeString,
		"ስህተት":  TypeError,
	}
}

func (m Amharic) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
