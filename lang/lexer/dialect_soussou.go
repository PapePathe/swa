package lexer

import "regexp"

type Soussou struct{}

var _ Dialect = (*Soussou)(nil)

func (m Soussou) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`konti`), defaultHandler(TypeInt, "konti")},
		{regexp.MustCompile(`sèbèli`), defaultHandler(TypeString, "sèbèli")},
		{regexp.MustCompile(`xamuara`), defaultHandler(KeywordElse, "xamuara")},
		{regexp.MustCompile(`xa`), defaultHandler(KeywordIf, "xa")},
		{regexp.MustCompile(`masen`), defaultHandler(Print, "masen")},
		{regexp.MustCompile(`structure`), defaultHandler(Struct, "structure")},
		{regexp.MustCompile(`labali`), defaultHandler(Main, "labali")},
		{regexp.MustCompile(`fonction`), defaultHandler(Function, "fonction")},
		{regexp.MustCompile(`ragbilen`), defaultHandler(Return, "ragbilen")},
		{
			regexp.MustCompile(
				`'[áàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z0-9]'`,
			),
			characterHandler,
		},
		{regexp.MustCompile(`\/\/.*`), commentHandler},
		{regexp.MustCompile(`"[^"]*"`), stringHandler},
		{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
		{
			regexp.MustCompile(
				`[a-zA-Z_][áàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z0-9_]*`,
			),
			symbolHandler,
		},
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

func (m Soussou) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"xa":        KeywordIf,
		"xamuara":   KeywordElse,
		"structure": Struct,
		"kouicé":    Let,
		"constante": Const,
		"konti":     TypeInt,
		"sèbèli":    TypeString,
	}
}
