package lexer

import (
	"fmt"
	"regexp"
	"swahili/lang/errmsg"
)

// Igbo
//
// Country: Nigeria
type Igbo struct{}

var _ Dialect = (*Igbo)(nil)

func (m Igbo) DetectionPattern() *regexp.Regexp {
	return regexp.MustCompile(`oluasusu:igbo;`)
}

func (m Igbo) Name() string {
	return "igbo"
}

func (m Igbo) Reserved() map[string]TokenKind {
	return map[string]TokenKind{
		"oluasusu":   DialectDeclaration,
		"ka":         Let,
		"ogidi":      Const,
		"maobu":      KeywordElse,
		"oburu":      KeywordIf,
		"dika":       KeywordWhile,
		"oru":        Function,
		"weghachi":   Return,
		"mbido":      Main,
		"deputa":     Print,
		"nhazi":      Struct,
		"ezigbo":     True,
		"ụgha":       False,
		"ezigha":     TypeBool,
		"mkpụrụ":     TypeByte,
		"ofe":        TypeFloat,
		"onuogugu":   TypeInt,
		"onuogugu64": TypeInt64,
		"ederede":    TypeString,
		"njehie":     TypeError,
		"ogologo":    Variadic,
		"efu":        Zero,
	}
}

func (m Igbo) Error(key string, args ...any) error {
	formatted, ok := m.translations()[key]

	if !ok {
		return fmt.Errorf("key %s does not exist in igbo dialect translations", key)
	}

	return errmsg.NewAstError(formatted, args...)
}

func (m Igbo) translations() map[string]string {
	return map[string]string{
		"ArrayAccessExpression.NameNotASymbol":             "Nkwupụta %v abụghị aha variable ziri ezi",
		"ArrayAccessExpression.NotFoundInSymbolTable":      "Variable %s adịghị na tebulu akara (symbol table)",
		"ArrayAccessExpression.AccessedIndexIsNotANumber":  "Naanị nọmba ka a na-anara dịka index array, nke dị ugbua: (%s)",
		"ArrayAccessExpression.NotFoundInArraySymbolTable": "Array (%s) adịghị na tebulu akara (symbol table)",
		"ArrayAccessExpression.IndexOutOfBounds":           "Ihe dị na index (%s) adịghị n'ime array (%s)",

		"NumberExpression.LessThanMinInt32":    "%d dị obere karịa ogo kacha nta maka int32",
		"NumberExpression.GreaterThanMaxInt32": "%d dị ukwuu karịa ogo kacha elu maka int32",

		"VisitStructDeclaration.SelfPointerReferenceNotAllowed": "A naghị anabata struct nwere pointer na-atụ aka na onwe ya, akụkụ: %s",
		"VisitStructDeclaration.SelfReferenceNotAllowed":        "A naghị anabata struct na-atụ aka na onwe ya, akụkụ: %s",
		"LLVMGenerator.VisitStructDeclaration.EmptyStruct":      "Struct (%s) ga-enwerịrị opekata mpe otu field",

		"VisitSymbolExpression.UnsupportedTypeAsGlobal": "Ụdị data %s a anaghị arụ ọrụ na global",
		"VisitPrefixExpression.OperatorNotSupported":    "PrefixExpression: operator %s anaghị arụ ọrụ",

		"LLVMTypeChecker.VisitVarDeclaration.UnexpectedValue":       "A tụrụ anya %s mana a hụrụ %s",
		"LLVMTypeChecker.VisitAssignmentExpression.UnexpectedValue": "A tụrụ anya nkenye nke %s mana a hụrụ %s",

		"CompilerCtx.UpdateStructSymbol.StructDoesNotExist": "Enweghị ike imelite struct akpọrọ %s n'ihi na ọ dịghị na tebulu akara",
		"CompilerCtx.AddStructSymbol.StructAlreadyExist":    "Struct akpọrọ %s adịlarị na tebulu akara",
		"CompilerCtx.FindStructSymbol.StructDoesNotExist":   "Struct akpọrọ %s adịghị na tebulu akara",
		"CompilerCtx.AddArraySymbol.AlreadyExisits":         "Array akpọrọ %s adịlarị na tebulu akara",
		"CompilerCtx.FindArraySymbol.DoesNotExist":          "Array akpọrọ %s adịghị na tebulu akara",
		"CompilerCtx.AddFuncSymbol.AlreadyExisits":          "Function akpọrọ %s adịlarị na tebulu akara",
		"CompilerCtx.AddSymbol.AlreadyExisits":              "Variable akpọrọ %s adịlarị na tebulu akara",
		"CompilerCtx.FindSymbol.DoesNotExist":               "Variable akpọrọ %s adịghị na tebulu akara",
		"CompilerCtx.FindFuncSymbol.DoesNotExist":           "Function akpọrọ %s adịghị na tebulu akara",

		"LLVMGenerator.VisitArrayAccessExpression.FieldDoesNotExistInStruct": "Struct %s enweghị field %s",
		"LLVMGenerator.VisitArrayAccessExpression.UnderlyingTypeNotSet":      "E nsetaghị underlying type",
		"LLVMGenerator.VisitArrayAccessExpression.MissingSymbolTableEntry":   "ArrayAccessExpression: SymbolTableEntry na-efu efu",
		"LLVMGenerator.VisitArrayAccessExpression.PropertyIsnotAnArray":      "Akụkụ %s abụghị array",
		"LLVMGenerator.VisitArrayAccessExpression.NotImplementedFor":         "Emebebeghị ArrayAccessExpression maka %T",

		"LLVMGenerator.VisitArrayInitializationExpression.NotInsideFunction":  "Mmalite array (array initialization) kwesịrị ime n'ime function",
		"LLVMGenerator.VisitArrayInitializationExpression.UnsupportedElement": "Ihe mmalite array anaghị arụ ọrụ: %T",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.NotImplementedFor": "Emebebeghị ArrayOfStructsAccessExpression maka %T",
		"LLVMGenerator.VisitArrayOfStructsAccessExpression.PropertyNotFound":  "Ahụghị akụkụ (%s) n'ime struct %s",

		"LLVMGenerator.VisitVarDeclaration.NotInsideFunction":          "Nkwupụta global var nwere %T anaghị arụ ọrụ",
		"LLVMGenerator.VisitVarDeclaration.AlreadyExisits":             "Variable %s adịlarị",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedInitializerType": "Nkwupụta var nwere %T anaghị arụ ọrụ",
		"LLVMGenerator.VisitVarDeclaration.UnsupportedTypeAsGlobal":    "Nkwupụta global var nwere %T anaghị arụ ọrụ",

		"LLVMGenerator.VisitBinaryExpression.StringsAreNotSupported": "A naghị anabata Strings na %s nke binary expression",
		"LLVMGenerator.VisitBinaryExpression.TypeMismatch":           "Ụdị data ekwekọghị na %v: %w",
		"LLVMGenerator.VisitBinaryExpression.CannotCoerceType":       "Enweghị ike ijikọta %s na %s",
		"LLVMGenerator.VisitBinaryExpression.UnsupportedOperator":    "Binary expressions: operator <%s> anaghị arụ ọrụ",

		"LLVMGenerator.VisitFunctionCall.UnsupportedType":                "Ụdị %T anaghị arụ ọrụ dịka argument na function call",
		"LLVMGenerator.VisitFunctionCall.DoesNotExist":                   "Function %s adịghị",
		"LLVMGenerator.VisitFunctionCall.ArgsAndParamsCountAreDifferent": "Function %s chọrọ argument %d mana e nyere ya %d",
		"LLVMGenerator.VisitFunctionCall.NameIsNotASymbol":               "FunctionCallExpression: aha abụghị akara (symbol)",
		"LLVMGenerator.VisitFunctionCall.FailedToEvaluate":               "A kụpụrụ n'ịtụle argument %d",
		"LLVMGenerator.VisitFunctionCall.StructPropertyValueTypeIsNil":   "E kwesịrị iseta struct property value type -- Ọzọkwa, enweghị ike inweta item array nke bụ struct",
		"LLVMGenerator.VisitFunctionCall.UnexpectedArgumentType":         "A tụrụ anya argument nke ụdị %s mana a hụrụ %s",

		"VisitReturnStatement.UnsupportedExpression":                    "VisitReturnStatement: Ụdị nkwupụta %T anaghị arụ ọrụ",
		"LLVMGenerator.VisitFunctionDefinition.UnsupportedArgumentType": "FuncDeclStatement argument ụdị %T anaghị arụ ọrụ",

		"LLVMGenerator.VisitMemberExpression.PropertyNotFound":   "Akụkụ akpọrọ %s adịghị na index %d",
		"LLVMGenerator.VisitMemberExpression.NotDefined":         "Variable %s adịghị",
		"LLVMGenerator.VisitMemberExpression.NotAStructInstance": "Variable %s abụghị struct instance",

		"LLVMGenerator.VisitStructInitializationExpression.Unimplemented":     "Emebebeghị struct field initialization maka %T",
		"LLVMGenerator.VisitStructInitializationExpression.NotInsideFunction": "Mmalite struct (struct initialization) kwesịrị ime n'ime function",

		"LLVMGenerator.resolveGepIndices.FailedToEvaluate": "A kụpụrụ n'ịtụle nkwupụta index (index expression)",

		"LLVMGenerator.getProperty.NotASymbol": "Struct property kwesịrị ịbụ akara (symbol)",

		"LLVMGenerator.ZeroOfArrayType.TooBigForZeroInitializer": "ArraySize (%d) buru ibu maka zero value initialization, nke kacha mma bụ %d",

		"LLVMCompiler.MissingProgramEntrypoint":  "Ihe omume gị (program) enweghị main function",
		"LLVMCompiler.TooManyProgramEntrypoints": "Ihe omume gị ga-enwerịrị naanị otu main function, ọnụọgụ ahụrụ (%d)",
	}
}
