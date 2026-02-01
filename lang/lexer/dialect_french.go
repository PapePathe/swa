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
		"constante":  Const,
		"dialecte":   DialectDeclaration,
		"fonction":   Function,
		"sinon":      KeywordElse,
		"si":         KeywordIf,
		"tantque":    KeywordWhile,
		"variable":   Let,
		"demarrer":   Main,
		"afficher":   Print,
		"retourner":  Return,
		"structure":  Struct,
		"decimal":    TypeFloat,
		"entier":     TypeInt,
		"entier64":   TypeInt64,
		"chaine":     TypeString,
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

		"VisitStructDeclaration.SelfPointerReferenceNotAllowed": "struct with pointer reference to self not supported, property: %s",
		"VisitStructDeclaration.SelfReferenceNotAllowed":        "struct with reference to self not supported, property: %s",

		"VisitSymbolExpression.UnsupportedTypeAsGlobal": "Unsupported datatype %s in global",
		"VisitPrefixExpression.OperatorNotSupported":    "PrefixExpression: operator %s not supported",

		"LLVMTypeChecker.VisitVarDeclaration.UnexpectedValue":       "expected %s but got %s",
		"LLVMTypeChecker.VisitAssignmentExpression.UnexpectedValue": "Expected assignment of %s but got %s",

		"CompilerCtx.UpdateStructSymbol.StructDoesNotExist": "struct named %s cannot be updated since it does not exist in symbol table",
		"CompilerCtx.AddStructSymbol.StructAlreadyExist":    "struct named %s already exists in symbol table",
		"CompilerCtx.FindStructSymbol.StructDoesNotExist":   "struct named %s does not exist in symbol table",
		"CompilerCtx.AddArraySymbol.AlreadyExisits":         "array named %s already exists in symbol table",
		"CompilerCtx.FindArraySymbol.DoesNotExist":          "array named %s does not exist in symbol table",
		"CompilerCtx.AddFuncSymbol.AlreadyExisits":          "function named %s already exists in symbol table",
		"CompilerCtx.AddSymbol.AlreadyExisits":              "variable named %s already exists in symbol table",
		"CompilerCtx.FindSymbol.DoesNotExist":               "variable named %s does not exist in symbol table",
	}
}
