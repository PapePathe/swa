package lexer

import "regexp"

type English struct{}

var _ Dialect = (*English)(nil)

func (m English) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`int`), defaultHandler(TypeInt, "int")},
		{regexp.MustCompile(`if`), defaultHandler(KeywordIf, "else")},
		{regexp.MustCompile(`else`), defaultHandler(KeywordElse, "if")},
		{regexp.MustCompile(`struct`), defaultHandler(Struct, "struct")},
		{regexp.MustCompile(`\/\/.*`), commentHandler},
		{regexp.MustCompile(`"[^"]*"`), stringHandler},
		{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
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

func (m English) Reserved() map[string]TokenKind {
	return map[string]TokenKind{}
}
