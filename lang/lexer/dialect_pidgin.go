package lexer

import (
	"regexp"
)

// Country: Nigeria
// Nigerian Pidgin (Naija) is the most widely spoken lingua franca in Nigeria, bridging the communication
// gap across its 250+ ethnic groups. It is an English-based creole that continues to evolve as a
// dynamic language of urban culture, music, and everyday interaction.
type Pidgin struct{}

var _ Dialect = (*Pidgin)(nil)

func (m Pidgin) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`pidgin:pidgin;`)
}

func (m Pidgin) Name() string {
	return "pidgin"
}

func (m Pidgin) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"mek":           Let,
		"no-fit-change": Const,
		"if":            KeywordIf,
		"else":          KeywordElse,
		"while":         KeywordWhile,
		"function":      Function,
		"return":        Return,
		"start":         Main,
		"show":          Print,
		"thing":         Struct,
		"na-so":         True,
		"no-be-so":      False,
		"bool":          TypeBool,
		"byte":          TypeByte,
		"decimal":       TypeFloat,
		"namba":         TypeInt,
		"namba64":       TypeInt64,
		"word":          TypeString,
		"wahala":        TypeError,
	}
}

func (m Pidgin) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
