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
		"erreur":     TypeError,
		"variadique": Variadic,
		"zero":       Zero,
	}
}

func (m French) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		format := "key %s does not exist in french dialect translations"

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
		"LLVMGenerator.VisitStructDeclaration.EmptyStruct":      "La structure (%s) doit avoir au moins un champ",

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
		"CompilerCtx.FindFuncSymbol.DoesNotExist":           "function named %s does not exist in symbol table",

		"LLVMGenerator.VisitArrayAccessExpression.FieldDoesNotExistInStruct": "struct %s has no field %s",
		"LLVMGenerator.VisitArrayAccessExpression.UnderlyingTypeNotSet":      "underlying type not set",
		"LLVMGenerator.VisitArrayAccessExpression.MissingSymbolTableEntry":   "ArrayAccessExpression Missing SymbolTableEntry",
		"LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray":      "ArrayAccessExpression property %s is not an array",
		"LLVMGenerator.VisitArrayAccessExpression.NotImplementedFor":         "ArrayAccessExpression not implemented for %T",

		"LLVMGenerator.VisitArrayInitializationExpression.NotInsideFunction":  "array initialization should happen inside a function",
		"LLVMGenerator.VisitArrayInitializationExpression.UnsupportedElement": "unsupported array initialization element: %T",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.NotImplementedFor": "ArrayOfStructsAccessExpression not implemented for %T",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.PropertyNotFound":  "property (%s) not found in struct %s",

		"LLVMGenerator.VisitVarDeclaration.NotInsideFunction":          "global var decl with %T not supported",
		"LLVMGenerator.VisitVarDeclaration.AlreadyExisits":             "variable %s is already defined",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedInitializerType": "var decl with %T not supported",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedTypeAsGlobal":    "global var decl with %T not supported",

		"LLVMGenerator.VisitBinaryExpression.StringsAreNotSupported": "Strings are not supported in %s of binary expression",
		"LLVMGenerator.VisitBinaryExpression.TypeMismatch":           "type mismatch at %v: %w",
		"LLVMGenerator.VisitBinaryExpression.CannotCoerceType":       "cannot coerce %s and %s",
		"LLVMGenerator.VisitBinaryExpression.UnsupportedOperator":    "Binary expressions : unsupported operator <%s>",

		"LLVMGenerator.VisitFunctionCall.UnsupportedType":                "Type %T not supported as function call argument",
		"LLVMGenerator.VisitFunctionCall.DoesNotExist":                   "function %s does not exist",
		"LLVMGenerator.VisitFunctionCall.ArgsAndParamsCountAreDifferent": "function %s expect %d arguments but was given %d",
		"LLVMGenerator.VisitFunctionCall.NameIsNotASymbol":               "FunctionCallExpression: name is not a symbol",
		"LLVMGenerator.VisitFunctionCall.FailedToEvaluate":               "failed to evaluate argument %d",
		"LLVMGenerator.VisitFunctionCall.StructPropertyValueTypeIsNil":   "Struct property value type should be set -- Also cannot index array item that is a struct",
		"LLVMGenerator.VisitFunctionCall.UnexpectedArgumentType":         "expected argument of type %s but got %s",

		"LLVMGenerator.VisitFunctionDefinition.UnsupportedArgumentType": "FuncDeclStatement argument type %v not supported",

		"LLVMGenerator.VisitMemberExpression.PropertyNotFound":   "Property named %s does not exist at index %d",
		"LLVMGenerator.VisitMemberExpression.NotDefined":         "variable %s is not defined",
		"LLVMGenerator.VisitMemberExpression.NotAStructInstance": "variable %s is not a struct instance",

		"LLVMGenerator.VisitStructInitializationExpression.Unimplemented": "struct field initialization unimplemented for %T",

		"LLVMGenerator.resolveGepIndices.FailedToEvaluate": "failed to evaluate index expression",
		"LLVMGenerator.getProperty.NotASymbol":             "struct property should be a symbol",

		"LLVMCompiler.MissingProgramEntrypoint":  "Vous devez definir le programme principal",
		"LLVMCompiler.TooManyProgramEntrypoints": "Vous devez definir un seul programme principal, nombre (%d)",
	}
}
