package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type French struct{}

var _ Dialect = (*French)(nil)

func (French) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`dialecte:français;`)
}

func (m French) Patterns() []RegexpPattern {
	return []RegexpPattern{
		{regexp.MustCompile(`\n`), newlineHandler},
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`\bdialecte\b`), defaultHandler(DialectDeclaration, "dialecte")},
		{regexp.MustCompile(`\bvariadique\b`), defaultHandler(Variadic, "variadique")},
		{regexp.MustCompile(`\bentier64\b`), defaultHandler(TypeInt64, "entier64")},
		{regexp.MustCompile(`\bentier\b`), defaultHandler(TypeInt, "entier")},
		{regexp.MustCompile(`\bdecimal\b`), defaultHandler(TypeFloat, "decimal")},
		{regexp.MustCompile(`\bchaine\b`), defaultHandler(TypeString, "chaine")},
		{regexp.MustCompile(`\bsinon\b`), defaultHandler(KeywordElse, "sinon")},
		{regexp.MustCompile(`\bsi\b`), defaultHandler(KeywordIf, "si")},
		{regexp.MustCompile(`tant que`), defaultHandler(KeywordWhile, "tant que")},
		{regexp.MustCompile(`\bafficher\b`), defaultHandler(Print, "afficher")},
		{regexp.MustCompile(`\bstructure\b`), defaultHandler(Struct, "structure")},
		{regexp.MustCompile(`\bdemarrer\b`), defaultHandler(Main, "demarrer")},
		{regexp.MustCompile(`\bfonction\b`), defaultHandler(Function, "fonction")},
		{regexp.MustCompile(`\bretourner\b`), defaultHandler(Return, "retourner")},
		{
			regexp.MustCompile(
				`\b[áàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z_][áàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽaa-zA-Z0-9_]*\b`,
			),
			symbolHandler,
		},
		{
			regexp.MustCompile(
				`'[áàâãäåæçćčđéèêëíìîïðñóòôõöøœśšşțúùûüýÿžAÁÀÂÃÄÅÆÇĆČĐÉÈÊËÍÌÎÏÐÑÓÒÔÕÖØŒŚŠŞȚÚÙÛÜÝŸŽa-zA-Z0-9]'`,
			),
			characterHandler,
		},
		{regexp.MustCompile(`\/\/.*`), commentHandler},
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
		{regexp.MustCompile(`:`), defaultHandler(Colon, ":")},
		{regexp.MustCompile(`,`), defaultHandler(Comma, ",")},
		{regexp.MustCompile(`\+`), defaultHandler(Plus, "+")},
		{regexp.MustCompile(`\%`), defaultHandler(Modulo, "%")},
		{regexp.MustCompile(`-`), defaultHandler(Minus, "-")},
		{regexp.MustCompile(`/`), defaultHandler(Divide, "/")},
		{regexp.MustCompile(`\*`), defaultHandler(Star, "*")},
	}
}

func (m French) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		format := "key %s does not exist in dialect translations"

		return fmt.Errorf(format, key)
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m French) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"afficher":   Print,
		"chaine":     TypeString,
		"constante":  Const,
		"decimal":    TypeFloat,
		"demarrer":   Main,
		"dialecte":   DialectDeclaration,
		"tantque":    KeywordWhile,
		"entier":     TypeInt,
		"fonction":   Function,
		"si":         KeywordIf,
		"sinon":      KeywordElse,
		"structure":  Struct,
		"variable":   Let,
		"variadique": Variadic,
	}
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
