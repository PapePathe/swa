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

func (m English) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\n`), newlineHandler},
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`\bdialect\b`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`\bint64\b`), defaultHandler(TypeInt64, "int64")},
		{regexp.MustCompile(`\bint\b`), defaultHandler(TypeInt, "int")},
		{regexp.MustCompile(`\bfloat\b`), defaultHandler(TypeFloat, "float")},
		{regexp.MustCompile(`\bstring\b`), defaultHandler(TypeString, "string")},
		{regexp.MustCompile(`\bwhile\b`), defaultHandler(KeywordWhile, "while")},
		{regexp.MustCompile(`\bif\b`), defaultHandler(KeywordIf, "if")},
		{regexp.MustCompile(`\belse\b`), defaultHandler(KeywordElse, "else")},
		{regexp.MustCompile(`\bprint\b`), defaultHandler(Print, "print")},
		{regexp.MustCompile(`\bstruct\b`), defaultHandler(Struct, "struct")},
		{regexp.MustCompile(`\bstart\b`), defaultHandler(Main, "start")},
		{regexp.MustCompile(`\breturn\b`), defaultHandler(Return, "return")},
		{regexp.MustCompile(`\bvariadic\b`), defaultHandler(Variadic, "variadic")},
		{regexp.MustCompile(`\bfunc\b`), defaultHandler(Function, "func")},
		{regexp.MustCompile(`[a-zA-Z_]([a-zA-Z0-9_])*`), symbolHandler},
		{regexp.MustCompile(`\/\/.*`), commentHandler},
		{regexp.MustCompile(`'[a-zA-Z0-9]'`), characterHandler},
		{regexp.MustCompile(`"[^"]*"`), stringHandler},
		{regexp.MustCompile(`[-]?[0-9]+\.[0-9]+`), floatHandler},
		{regexp.MustCompile(`[-]?[0-9]+`), numberHandler},
		{regexp.MustCompile(`\[`), defaultHandler(OpenBracket, "[")},
		{regexp.MustCompile(`\]`), defaultHandler(CloseBracket, "]")},
		{regexp.MustCompile(`\{`), defaultHandler(OpenCurly, "{")},
		{regexp.MustCompile(`\}`), defaultHandler(CloseCurly, "}")},
		{regexp.MustCompile(`\(`), defaultHandler(OpenParen, "(")},
		{regexp.MustCompile(`\)`), defaultHandler(CloseParen, ")")},
		{regexp.MustCompile(`!=`), defaultHandler(NotEquals, "!=")},
		{regexp.MustCompile(`\+=`), defaultHandler(PlusEquals, "+=")},
		{regexp.MustCompile(`==`), defaultHandler(Equals, "==")},
		{regexp.MustCompile(`=`), defaultHandler(Assignment, "=")},
		{regexp.MustCompile(`!`), defaultHandler(Not, "!")},
		{regexp.MustCompile(`<=`), defaultHandler(LessThanEquals, "<=")},
		{regexp.MustCompile(`<`), defaultHandler(LessThan, "<")},
		{regexp.MustCompile(`>=`), defaultHandler(GreaterThanEquals, ">=")},
		{regexp.MustCompile(`>`), defaultHandler(GreaterThan, ">")},
		{regexp.MustCompile(`\|\|`), defaultHandler(Or, "||")},
		{regexp.MustCompile(`&&`), defaultHandler(And, "&&")},
		{regexp.MustCompile(`\.`), defaultHandler(Dot, ".")},
		{regexp.MustCompile(`;`), defaultHandler(SemiColon, ";")},
		{regexp.MustCompile(`::`), defaultHandler(DoubleColon, "::")},
		{regexp.MustCompile(`:`), defaultHandler(Colon, ":")},
		{regexp.MustCompile(`,`), defaultHandler(Comma, ",")},
		{regexp.MustCompile(`\+`), defaultHandler(Plus, "+")},
		{regexp.MustCompile(`\%`), defaultHandler(Modulo, "%")},
		{regexp.MustCompile(`-`), defaultHandler(Minus, "-")},
		{regexp.MustCompile(`/`), defaultHandler(Divide, "/")},
		{regexp.MustCompile(`\*`), defaultHandler(Star, "*")},
	}
}

func (m English) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		panic(fmt.Sprintf("key %s does not exist in dialect translations", key))
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m English) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"if":     KeywordIf,
		"else":   KeywordElse,
		"struct": Struct,
		"func":   Function,
		"let":    Let,
		"const":  Const,
		"int":    TypeInt,
		"float":  TypeFloat,
		"string": TypeString,
	}
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
