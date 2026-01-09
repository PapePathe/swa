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

func (m French) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"afficher":   Print,
		"chaine":     TypeString,
		"constante":  Const,
		"decimal":    TypeFloat,
		"demarrer":   Main,
		"dialecte":   DialectDeclaration,
		"entier":     TypeInt,
		"entier64":   TypeInt64,
		"fonction":   Function,
		"retourner":  Return,
		"si":         KeywordIf,
		"sinon":      KeywordElse,
		"structure":  Struct,
		"tantque":    KeywordWhile,
		"variable":   Let,
		"variadique": Variadic,
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

func (m French) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "L'expression (%s) n'est pas un nom de variable correct",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "La variable %s n'existe pas dans la table des symboles",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Seuls les nombres positif sont permis comme indice de tableau, valeur courante: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Le tableau (%s) n'existe pas dans la table des symboles",
		"ArrayAccessExpression.IndexOutOfBounds":           "L'element a la position (%s) depasse les limites du tableau (%s)",

		"NumberExpression.LessThanMinInt32":    "%d est plus petit que la valeur minimale pour un entier de 32 bits",
		"NumberExpression.GreaterThanMaxInt32": "%d est plus grand que la valeur maximale pour un entier de 32 bits",
	}
}
