package lexer

import (
	"fmt"
	"regexp"
	"swahili/lang/errmsg"
)

type Soussou struct{}

var _ Dialect = (*Soussou)(nil)

func (Soussou) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialect:soussou;`)
}

func (m Soussou) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`konti`), defaultHandler(TypeInt, "konti")},
		{regexp.MustCompile(`sèbèli`), defaultHandler(TypeString, "sèbèli")},
		{regexp.MustCompile(`xamuara`), defaultHandler(KeywordElse, "xamuara")},
		{regexp.MustCompile(`xa`), defaultHandler(KeywordIf, "xa")},
		{regexp.MustCompile(`masen`), defaultHandler(Print, "masen")},
		{regexp.MustCompile(`fokhi`), defaultHandler(Struct, "fokhi")},
		{regexp.MustCompile(`sodé`), defaultHandler(Main, "sodé")},
		{regexp.MustCompile(`wali`), defaultHandler(Function, "wali")},
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
		"fokhi":     Struct,
		"kouicé":    Let,
		"constante": Const,
		"konti":     TypeInt,
		"sèbèli":    TypeString,
		"sodé":      Main,
		"wali":      Function,
		"masen":     Print,
	}
}

func (m Soussou) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		panic(fmt.Sprintf("key %s does not exist in dialect translations", key))
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m Soussou) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "L'expression (%s) n'est pas un nom de variable correct",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "La variable %s n'existe pas dans la table des symboles",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Le tableau (%s) n'existe pas dans la table des symboles",
		"ArrayAccessExpression.IndexOutOfBounds":           "L'element a la position (%s) depasse les limites du tableau (%s)",
	}
}
