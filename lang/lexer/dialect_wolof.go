package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Wolof struct{}

var _ Dialect = (*Wolof)(nil)

func (m Wolof) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`laak:wolof;`)
}

func (m Wolof) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in wolof dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args)
}

func (m Wolof) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"constante": Const,
		"laak":      DialectDeclaration,
		// TODO: find  translation for function
		"fonction": Function,
		"sude":     KeywordIf,
		"wala":     KeywordElse,
		// TODO: find  translation for while
		"while":     KeywordWhile,
		"dencukaay": Let,
		"tambali":   Main,
		"wanel":     Print,
		"deloko":    Return,
		"structure": Struct,
		"decimal":   TypeFloat,
		"lëmm":      TypeInt,
		"lëmm64":    TypeInt64,
		"ay_araf":   TypeString,
		"njumte":    TypeError,
		// TODO: find  translation for variadic
		"variadique": Variadic,
		"zero":       Zero,
	}
}

func (m Wolof) translations() map[string]string {
	// TODO: add translations for error messages
	return map[string]string{}
}
