package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Malinke struct{}

var _ Dialect = (*Malinke)(nil)

func (Malinke) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile("kan:malinke;")
}

func (Malinke) Name() string {
	return "malinke"
}

func (m Malinke) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in malinke dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m Malinke) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"tintin":    Const,
		"kan":       DialectDeclaration,
		"baara":     Function,
		"wala":      KeywordElse,
		"ni":        KeywordIf,
		"tuma":      KeywordWhile,
		"atö":       Let,
		"daminen":   Main,
		"yira":      Print,
		"segin":     Return,
		"joyoro":    Struct,
		"jatelen":   TypeFloat,
		"jate":      TypeInt,
		"jate64":    TypeInt64,
		"seben":     TypeString,
		"fili":      TypeError,
		"caaman":    Variadic,
		"foy":       Zero,
		"tinye":     True,
		"wouya":     False,
		"tinyejate": TypeBool,
		"make":      Make,
		"count":     Len,
		"cap":       Cap,
		"append":    Append,
	}
}

func (m Malinke) translations() map[string]string {
	return map[string]string{
		"AppendExpression.AppendableIsNotASlice": "append expects a slice, got %T",
		"ListExpression.MakeRequiresSliceType":   "make() function requires a slice allocation map, got %v",
		"ListCountExpression.RequiresSliceType":  "count() function requires a slice, got %v",

		"LLVMCompiler.MissingProgramEntrypoint":  "Your program is missing a main function",
		"LLVMCompiler.TooManyProgramEntrypoints": "Your program must have exactly one main function, count (%d)",
	}
}
