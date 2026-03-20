package lexer

import (
	"regexp"
)

// Country: Togo
// Ikposo (Kposo) is a Kwa language spoken by the Akposso people, primarily in
// the Plateaux Region of Togo and parts of Ghana. It is the major language
// of the Akposso people and has a unique linguistic heritage.
type Ikposo struct{}

var _ Dialect = (*Ikposo)(nil)

func (m Ikposo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`ikposo:ikposo;`)
}

func (m Ikposo) Name() string {
	return "ikposo"
}

func (m Ikposo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"na":      Let,
		"ma-trɔ":  Const,
		"ne":      KeywordIf,
		"ne-ba":   KeywordElse,
		"esi":     KeywordWhile,
		"dɔ":      Function,
		"trɔ":     Return,
		"gɔme":    Main,
		"ŋlɔ":     Print,
		"ɖoɖo":    Struct,
		"nyateƒe": True,
		"alakpa":  False,
		"bool":    TypeBool,
		"byte":    TypeByte,
		"desimal": TypeFloat,
		"lamba":   TypeInt,
		"lamba64": TypeInt64,
		"nya":     TypeString,
		"vodada":  TypeError,
	}
}

func (m Ikposo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
