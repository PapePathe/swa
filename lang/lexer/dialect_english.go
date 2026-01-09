package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type English struct{}

var _ Dialect = (*English)(nil)

func (m English) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:english;`)
}

func (m English) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"const":    Const,
		"dialect":  DialectDeclaration,
		"else":     KeywordElse,
		"float":    TypeFloat,
		"func":     Function,
		"if":       KeywordIf,
		"int":      TypeInt,
		"int64":    TypeInt64,
		"let":      Let,
		"print":    Print,
		"return":   Return,
		"string":   TypeString,
		"start":    Main,
		"struct":   Struct,
		"variadic": Variadic,
		"while":    KeywordWhile,
	}
}

func (m English) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in english dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m English) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "The expression %v is not a correct variable name",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "The variable %s does not exist in symbol table",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Only numbers are supported as array index, current: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Array (%s) does not exist in symbol table",
		"ArrayAccessExpression.IndexOutOfBounds":           "Element at index (%s) does not exist in array (%s)",
	}
}
