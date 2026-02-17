package ast

import "swahili/lang/lexer"

// Type
type Type interface {
	Accept(g CodeGenerator) error
	AcceptZero(g CodeGenerator) error
	Value() DataType
	Equals(Type) bool
	String() string
}

type ExpressionType interface {
	Equivalent(a ExpressionType, b ExpressionType) bool
	Compatible(a ExpressionType, b ExpressionType) bool
}

// Node ...
type Node interface {
	Accept(g CodeGenerator) error
	TokenStream() []lexer.Token
}

// Statement ...
type Statement interface {
	Node
}

// Expression ...
type Expression interface {
	Node
	VisitedSwaType() Type
}

// CodeGenerator ...
type CodeGenerator interface {
	// Statements
	StatementsCodeGenerator
	// Expressions
	ExpressionsCodeGenerator
	// Types
	TypeVisitor
	// Zero values
	ZeroValueVisitor
}

type StatementsCodeGenerator interface {
	VisitBlockStatement(node *BlockStatement) error
	VisitVarDeclaration(node *VarDeclarationStatement) error
	VisitReturnStatement(node *ReturnStatement) error
	VisitExpressionStatement(node *ExpressionStatement) error
	VisitConditionalStatement(node *ConditionalStatetement) error
	VisitWhileStatement(node *WhileStatement) error
	VisitFunctionDefinition(node *FuncDeclStatement) error
	VisitStructDeclaration(node *StructDeclarationStatement) error
	VisitPrintStatement(node *PrintStatetement) error
	VisitMainStatement(node *MainStatement) error
}

type ExpressionsCodeGenerator interface {
	VisitBinaryExpression(node *BinaryExpression) error
	VisitFunctionCall(node *FunctionCallExpression) error
	VisitCallExpression(node *CallExpression) error
	VisitStringExpression(node *StringExpression) error
	VisitNumberExpression(node *NumberExpression) error
	VisitFloatExpression(node *FloatExpression) error
	VisitSymbolExpression(node *SymbolExpression) error
	VisitAssignmentExpression(node *AssignmentExpression) error
	VisitTupleExpression(node *TupleExpression) error
	VisitTupleAssignmentExpression(node *TupleAssignmentExpression) error
	VisitErrorExpression(node *ErrorExpression) error
	VisitMemberExpression(node *MemberExpression) error
	VisitArrayAccessExpression(node *ArrayAccessExpression) error
	VisitStructInitializationExpression(node *StructInitializationExpression) error
	VisitArrayInitializationExpression(node *ArrayInitializationExpression) error
	VisitArrayOfStructsAccessExpression(node *ArrayOfStructsAccessExpression) error
	VisitPrefixExpression(node *PrefixExpression) error
	VisitZeroExpression(node *ZeroExpression) error
	VisitFloatingBlockExpression(node *FloatingBlockExpression) error
	VisitSymbolValueExpression(node *SymbolValueExpression) error
	VisitSymbolAdressExpression(node *SymbolAdressExpression) error
	VisitBooleanExpression(node *BooleanExpression) error
}

type TypeVisitor interface {
	VisitSymbolType(node *SymbolType) error
	VisitTupleType(node *TupleType) error
	VisitNumberType(node *NumberType) error
	VisitNumber64Type(node *Number64Type) error
	VisitFloatType(node *FloatType) error
	VisitErrorType(node *ErrorType) error
	VisitPointerType(node *PointerType) error
	VisitStringType(node *StringType) error
	VisitArrayType(node *ArrayType) error
	VisitVoidType(node *VoidType) error
	VisitBoolType(node *BoolType) error
}

type ZeroValueVisitor interface {
	ZeroOfSymbolType(node *SymbolType) error
	ZeroOfTupleType(node *TupleType) error
	ZeroOfNumberType(node *NumberType) error
	ZeroOfNumber64Type(node *Number64Type) error
	ZeroOfFloatType(node *FloatType) error
	ZeroOfErrorType(node *ErrorType) error
	ZeroOfPointerType(node *PointerType) error
	ZeroOfStringType(node *StringType) error
	ZeroOfArrayType(node *ArrayType) error
	ZeroOfVoidType(node *VoidType) error
	ZeroOfBoolType(node *BoolType) error
}
