package lexer

import (
	"fmt"
	"regexp"
	"swahili/lang/errmsg"
)

type French struct{}

var _ Dialect = (*French)(nil)

func (m French) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`dialect`), defaultHandler(DialectDeclaration, "dialect")},
		{regexp.MustCompile(`entier`), defaultHandler(TypeInt, "entier")},
		{regexp.MustCompile(`chaine`), defaultHandler(TypeString, "chaine")},
		{regexp.MustCompile(`sinon`), defaultHandler(KeywordElse, "sinon")},
		{regexp.MustCompile(`si`), defaultHandler(KeywordIf, "si")},
		{regexp.MustCompile(`afficher`), defaultHandler(Print, "afficher")},
		{regexp.MustCompile(`structure`), defaultHandler(Struct, "structure")},
		{regexp.MustCompile(`demarrer`), defaultHandler(Main, "demarrer")},
		{regexp.MustCompile(`fonction`), defaultHandler(Function, "fonction")},
		{regexp.MustCompile(`retourner`), defaultHandler(Return, "retourner")},
		{
			regexp.MustCompile(
				`'[aáàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z0-9]'`,
			),
			characterHandler,
		},
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

func (m French) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		panic(fmt.Sprintf("key %s does not exist in dialect translations", key))
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m French) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "L'expression (%s) n'est pas un nom de variable correct",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "La variable %s n'existe pas dans la table des symboles",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Le tableau (%s) n'existe pas dans la table des symboles",
		"ArrayAccessExpression.IndexOutOfBounds":           "L'element a la position (%s) depasse les limites du tableau (%s)",
	}
}

func (m French) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"si":        KeywordIf,
		"sinon":     KeywordElse,
		"structure": Struct,
		"variable":  Let,
		"constante": Const,
		"entier":    TypeInt,
		"chaine":    TypeString,
	}
}
