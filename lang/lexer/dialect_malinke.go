package lexer

import (
	"fmt"
	"regexp"
)

type Malinke struct{}

var _ Dialect = (*Malinke)(nil)

func (Malinke) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:malinke;`)
}
func (m Malinke) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\n`), newlineHandler},
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`fèndo`), defaultHandler(TypeInt, "fèndo")},
		{regexp.MustCompile(`nii`), defaultHandler(KeywordElse, "nii")},
		{regexp.MustCompile(`ni`), defaultHandler(KeywordIf, "ni")},
		{regexp.MustCompile(`afo`), defaultHandler(Print, "afo")},
		{regexp.MustCompile(`struct`), defaultHandler(Struct, "struct")},
		{
			regexp.MustCompile(
				`'[aáàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z0-9]'`,
			),
			characterHandler,
		},
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

func (m Malinke) Error(key string, args ...any) error {
	return fmt.Errorf("test")
}

func (m Malinke) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"ni":     KeywordIf,
		"nii":    KeywordElse,
		"struct": Struct,
		"let":    Let,
		"const":  Const,
		"fèndo":  TypeInt,
	}
}
