package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Wolof struct{}

var _ Dialect = (*Wolof)(nil)

func (m Wolof) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:wolof;`)
}

func (m Wolof) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\n`), newlineHandler},
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`sifaru lëmm`), defaultHandler(TypeInt, "sifaru lëmm")},
		{regexp.MustCompile(`taxawlu araf`), defaultHandler(TypeString, "taxawlu araf")},
		{regexp.MustCompile(`wala`), defaultHandler(KeywordElse, "wala")},
		{regexp.MustCompile(`sude`), defaultHandler(KeywordIf, "sude")},
		{regexp.MustCompile(`struct`), defaultHandler(Struct, "struct")},
		{regexp.MustCompile(`tambali`), defaultHandler(Main, "tambali")},
		{regexp.MustCompile(`fonction`), defaultHandler(Function, "fonction")},
		{regexp.MustCompile(`deloko`), defaultHandler(Return, "deloko")},
		{regexp.MustCompile(`bindal`), defaultHandler(Print, "bindal")},
		{regexp.MustCompile(`\/\/.*`), commentHandler},
		{regexp.MustCompile(`"[^"]*"`), stringHandler},
		{regexp.MustCompile(`[-]?[0-9]+\.[0-9]+`), floatHandler},
		{regexp.MustCompile(`[-]?[0-9]+`), numberHandler},
		{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
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
		{regexp.MustCompile(`:`), defaultHandler(Colon, ":")},
		{regexp.MustCompile(`,`), defaultHandler(Comma, ",")},
		{regexp.MustCompile(`\+`), defaultHandler(Plus, "+")},
		{regexp.MustCompile(`-`), defaultHandler(Minus, "-")},
		{regexp.MustCompile(`/`), defaultHandler(Divide, "/")},
		{regexp.MustCompile(`\*`), defaultHandler(Star, "*")},
	}
}

func (m Wolof) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		panic(fmt.Sprintf("key %s does not exist in dialect translations", key))
	}

	return errmsg.NewAstError(formatted, args)
}

func (m Wolof) translations() map[string]string {
	// TODO: add reserved words
	return map[string]string{}
}

func (m Wolof) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"sude":         KeywordIf,
		"deloko":       Return,
		"tambali":      Main,
		"wala":         KeywordElse,
		"structure":    Struct,
		"dencukaay":    Let,
		"constante":    Const,
		"sifaru lëmm":  TypeInt,
		"taxawlu araf": TypeString,
	}
}
