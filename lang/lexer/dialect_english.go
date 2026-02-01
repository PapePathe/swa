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
		"func":     Function,
		"else":     KeywordElse,
		"if":       KeywordIf,
		"while":    KeywordWhile,
		"let":      Let,
		"start":    Main,
		"print":    Print,
		"return":   Return,
		"struct":   Struct,
		"float":    TypeFloat,
		"int":      TypeInt,
		"int64":    TypeInt64,
		"string":   TypeString,
		"variadic": Variadic,
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
		"ArrayAccessExpression.NameNotASymbol":                      "The expression %v is not a correct variable name",
		"ArrayAccessExpression.NotFoundInSymbolTable":               "The variable %s does not exist in symbol table",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":           "Only numbers are supported as array index, current: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable":          "Array (%s) does not exist in symbol table",
		"ArrayAccessExpression.IndexOutOfBounds":                    "Element at index (%s) does not exist in array (%s)",
		"NumberExpression.LessThanMinInt32":                         "%d is smaller than min value for int32",
		"NumberExpression.GreaterThanMaxInt32":                      "%d is greater than max value for int32",
		"VisitStructDeclaration.SelfPointerReferenceNotAllowed":     "struct with pointer reference to self not supported, property: %s",
		"VisitStructDeclaration.SelfReferenceNotAllowed":            "struct with reference to self not supported, property: %s",
		"VisitSymbolExpression.UnsupportedTypeAsGlobal":             "Unsupported datatype %s in global",
		"VisitPrefixExpression.OperatorNotSupported":                "PrefixExpression: operator %s not supported",
		"LLVMTypeChecker.VisitVarDeclaration.UnexpectedValue":       "expected %s but got %s",
		"LLVMTypeChecker.VisitAssignmentExpression.UnexpectedValue": "Expected assignment of %s but got %s",
	}
}
