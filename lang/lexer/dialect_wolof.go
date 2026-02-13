package lexer

import (
	"fmt"
	"regexp"

	"swahili/lang/errmsg"
)

type Wolof struct{}

var _ Dialect = (*Wolof)(nil)

func (m Wolof) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`laak:wolof;`)
}

func (m Wolof) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in wolof dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args)
}

func (m Wolof) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"constante": Const,
		"laak":      DialectDeclaration,
		// TODO: find  translation for function
		"fonction": Function,
		"sude":     KeywordIf,
		"wala":     KeywordElse,
		// TODO: find  translation for while
		"while":     KeywordWhile,
		"dencukaay": Let,
		"tambali":   Main,
		"wanel":     Print,
		"deloko":    Return,
		"structure": Struct,
		"decimal":   TypeFloat,
		"lëmm":      TypeInt,
		"lëmm64":    TypeInt64,
		"ay_araf":   TypeString,
		"njumte":    TypeError,
		// TODO: find  translation for variadic
		"variadique": Variadic,
		"zero":       Zero,
	}
}

func (m Wolof) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "The expression %v is not a correct variable name",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "The variable %s does not exist in symbol table",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Only numbers are supported as array index, current: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Array (%s) does not exist in symbol table",
		"ArrayAccessExpression.IndexOutOfBounds":           "Element at index (%s) does not exist in array (%s)",

		"NumberExpression.LessThanMinInt32":    "%d is smaller than min value for int32",
		"NumberExpression.GreaterThanMaxInt32": "%d is greater than max value for int32",

		"VisitStructDeclaration.SelfPointerReferenceNotAllowed": "struct with pointer reference to self not supported, property: %s",
		"VisitStructDeclaration.SelfReferenceNotAllowed":        "struct with reference to self not supported, property: %s",
		"LLVMGenerator.VisitStructDeclaration.EmptyStruct":      "Struct (%s) must have at least one field",

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

		"LLVMGenerator.VisitArrayAccessExpression.FieldDoesNotExistInStruct": "struct %s has no field %s",
		"LLVMGenerator.VisitArrayAccessExpression.UnderlyingTypeNotSet":      "underlying type not set",
		"LLVMGenerator.VisitArrayAccessExpression.MissingSymbolTableEntry":   "ArrayAccessExpression Missing SymbolTableEntry",
		"LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray":      "Property %s is not an array",
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

		"VisitReturnStatement.UnsupportedExpression":                    "VisitReturnStatement: Unsupported expression type %T",
		"LLVMGenerator.VisitFunctionDefinition.UnsupportedArgumentType": "FuncDeclStatement argument type %T not supported",

		"LLVMGenerator.VisitMemberExpression.PropertyNotFound":   "Property named %s does not exist at index %d",
		"LLVMGenerator.VisitMemberExpression.NotDefined":         "variable %s is not defined",
		"LLVMGenerator.VisitMemberExpression.NotAStructInstance": "variable %s is not a struct instance",

		"LLVMGenerator.VisitStructInitializationExpression.Unimplemented":     "struct field initialization unimplemented for %T",
		"LLVMGenerator.VisitStructInitializationExpression.NotInsideFunction": "struct initialization should happen inside a function",

		"LLVMGenerator.resolveGepIndices.FailedToEvaluate": "failed to evaluate index expression",

		"LLVMGenerator.getProperty.NotASymbol": "struct property should be a symbol",

		"LLVMGenerator.ZeroOfArrayType.TooBigForZeroInitializer": "ArraySize (%d) too big for zero value initialization, max is %d",
	}
}
