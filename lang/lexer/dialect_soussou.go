package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Soussou struct{}

var _ Dialect = (*Soussou)(nil)

func (m Soussou) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`khuien:soussou;`)
}

func (m Soussou) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"danyi":    Const,
		"khuien":   DialectDeclaration,
		"wali":     Function,
		"xamuara":  KeywordElse,
		"xa":       KeywordIf,
		"be":       KeywordWhile,
		"kouicé":   Let,
		"sodé":     Main,
		"masen":    Print,
		"gbilen":   Return,
		"fokhi":    Struct,
		"desimali": TypeFloat,
		"konti":    TypeInt,
		"kont64":   TypeInt64,
		"sèbèli":   TypeString,
		"foutoura": TypeError,
		"gbegbe":   Variadic,
		"maba":     Zero,
	}
}

func (m Soussou) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		format := "key %s does not exist in dialect translations"

		return fmt.Errorf(format, key)
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m Soussou) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "Kuitunyi (%s) mu variable xili ra",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "Variable %s maboroma",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Konti nan lanné tableau kui, yakosi: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Tableau (%s) maboroma",
		"ArrayAccessExpression.IndexOutOfBounds":           "Se kui (%s) mu na tableau (%s) kui",

		"NumberExpression.LessThanMinInt32":    "%d mu lanyi konti 32 bit xungbe ma",
		"NumberExpression.GreaterThanMaxInt32": "%d danyi konti 32 bit xungbe ma",

		"VisitStructDeclaration.SelfPointerReferenceNotAllowed": "fokhi ya na a gbe gbegbe findi, se mu lanné: %s",
		"VisitStructDeclaration.SelfReferenceNotAllowed":        "fokhi ya na a gbe findi, se mu lanné: %s",

		"VisitSymbolExpression.UnsupportedTypeAsGlobal": "DataType %s mu lanné global ra",
		"VisitPrefixExpression.OperatorNotSupported":    "Operator %s mu lanné prefix ma",

		"LLVMTypeChecker.VisitVarDeclaration.UnexpectedValue":       "%s nan lanné, kono %s nan masen",
		"LLVMTypeChecker.VisitAssignmentExpression.UnexpectedValue": "%s nan lanné %s binyé ma, kono %s nan masen",

		"CompilerCtx.UpdateStructSymbol.StructDoesNotExist": "fokhi %s maboroma, a mu noma findide",
		"CompilerCtx.AddStructSymbol.StructAlreadyExist":    "fokhi %s na na yi khorun",
		"CompilerCtx.FindStructSymbol.StructDoesNotExist":   "fokhi %s maboroma",
		"CompilerCtx.AddArraySymbol.AlreadyExisits":         "tableau %s na na yi khorun",
		"CompilerCtx.FindArraySymbol.DoesNotExist":          "tableau %s maboroma",
		"CompilerCtx.AddFuncSymbol.AlreadyExisits":          "wali %s na na yi khorun",
		"CompilerCtx.AddSymbol.AlreadyExisits":              "variable %s na na yi khorun",
		"CompilerCtx.FindSymbol.DoesNotExist":               "variable %s maboroma",

		"LLVMGenerator.VisitArrayAccessExpression.FieldDoesNotExistInStruct": "fokhi %s, se %s mu na a kui",
		"LLVMGenerator.VisitArrayAccessExpression.UnderlyingTypeNotSet":      "underlying type mu lanyi",
		"LLVMGenerator.VisitArrayAccessExpression.MissingSymbolTableEntry":   "ArrayAccessExpression SymbolTableEntry maboroma",
		"LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray":      "Property %s mu tableau ra",
		"LLVMGenerator.VisitArrayAccessExpression.NotImplementedFor":         "ArrayAccessExpression mu wali yi fe %T ra",

		"LLVMGenerator.VisitArrayInitializationExpression.NotInsideFunction":  "tableau lanyinyi wali kui nan lanyi",
		"LLVMGenerator.VisitArrayInitializationExpression.UnsupportedElement": "tableau element %T mu lanné",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.NotImplementedFor": "ArrayOfStructsAccessExpression mu wali yi fe %T ra",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.PropertyNotFound":  "se (%s) mu na fokhi %s kui",

		"LLVMGenerator.VisitVarDeclaration.NotInsideFunction":          "global variable %T mu lanné",
		"LLVMGenerator.VisitVarDeclaration.AlreadyExisits":             "variable %s na na yi khorun",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedInitializerType": "var decl %T mu lanné",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedTypeAsGlobal":    "global var decl %T mu lanné",

		"LLVMGenerator.VisitBinaryExpression.StringsAreNotSupported": "Sèbèli mu lanné %s binary expression kui",
		"LLVMGenerator.VisitBinaryExpression.TypeMismatch":           "type mismatch at %v: %w",
		"LLVMGenerator.VisitBinaryExpression.CannotCoerceType":       "%s nun %s mu noma lanyide",
		"LLVMGenerator.VisitBinaryExpression.UnsupportedOperator":    "Binary expression: operator <%s> mu lanné",

		"LLVMGenerator.VisitFunctionCall.UnsupportedType":                "Type %T mu lanné wali argument ra",
		"LLVMGenerator.VisitFunctionCall.DoesNotExist":                   "wali %s maboroma",
		"LLVMGenerator.VisitFunctionCall.ArgsAndParamsCountAreDifferent": "wali %s expect %d arguments kono %d nan masen",
		"LLVMGenerator.VisitFunctionCall.NameIsNotASymbol":               "FunctionCallExpression: xili mu variable ra",
		"LLVMGenerator.VisitFunctionCall.FailedToEvaluate":               "noma argument %d maboroma",
		"LLVMGenerator.VisitFunctionCall.StructPropertyValueTypeIsNil":   "fokhi se type mu lanyi -- also mu noma tableau se index ra fokhi kui",
		"LLVMGenerator.VisitFunctionCall.UnexpectedArgumentType":         "wali expect %s kono %s nan masen",

		"LLVMGenerator.VisitFunctionDefinition.UnsupportedArgumentType": "wali argument type %v mu lanné",

		"LLVMGenerator.VisitMemberExpression.PropertyNotFound":   "se %s mu lanyi index %d ra",
		"LLVMGenerator.VisitMemberExpression.NotDefined":         "variable %s maboroma",
		"LLVMGenerator.VisitMemberExpression.NotAStructInstance": "variable %s mu fokhi ra",

		"LLVMGenerator.VisitStructInitializationExpression.Unimplemented": "fokhi field initialization mu lanyi fe %T ra",

		"LLVMGenerator.resolveGepIndices.FailedToEvaluate": "index expression maboroma",

		"LLVMGenerator.getProperty.NotASymbol": "fokhi se nan lanné variable ra",
	}
}
