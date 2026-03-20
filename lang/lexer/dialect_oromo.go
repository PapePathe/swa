package lexer

import (
	"regexp"
)

// Country: Ethiopia
// Oromo is a Cushitic language spoken primarily by the Oromo people in
// Ethiopia and parts of Kenya. It is the most widely spoken language in
// Ethiopia and is central to the cultural and social life of the Oromia region.
type Oromo struct{}

var _ Dialect = (*Oromo)(nil)

func (m Oromo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`oromo:oromo;`)
}

func (m Oromo) Name() string {
	return "oromo"
}

func (m Oromo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"haa_tahu":        Let,
		"dhaabbataa":      Const,
		"yoo":             KeywordIf,
		"yoo_ta'uu_baate": KeywordElse,
		"yeroo":           KeywordWhile,
		"dalagaa":         Function,
		"deebisi":         Return,
		"ijoo":            Main,
		"maxxansi":        Print,
		"caasaa":          Struct,
		"dhugaa":          True,
		"soba":            False,
		"bool":            TypeBool,
		"byte":            TypeByte,
		"float":           TypeFloat,
		"lakkoofsa":       TypeInt,
		"lakkoofsa64":     TypeInt64,
		"jecha":           TypeString,
		"dogoggora":       TypeError,
	}
}

func (m Oromo) Error(key string, args ...any) error {
	return English{}.Error(key, args...)
}
