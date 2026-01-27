package compiler

import "swahili/lang/ast"

type LLVMTypeChecker struct {
	ctx *CompilerCtx
}

func NewLLVMTypeChecker(ctx *CompilerCtx) *LLVMTypeChecker {
	return &LLVMTypeChecker{ctx: ctx}
}

// VisitArrayAccessExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	return nil
}

// VisitArrayInitializationExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	return nil
}

// VisitArrayOfStructsAccessExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	return nil
}

// VisitArrayType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitArrayType(node *ast.ArrayType) error {
	return nil
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	return nil
}

// VisitBinaryExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitBinaryExpression(node *ast.BinaryExpression) error {
	return nil
}

// VisitBlockStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitBlockStatement(node *ast.BlockStatement) error {
	return nil
}

// VisitCallExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitCallExpression(node *ast.CallExpression) error {
	return nil
}

// VisitConditionalStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	return nil
}

// VisitExpressionStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return nil
}

// VisitFloatExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitFloatExpression(node *ast.FloatExpression) error {
	return nil
}

// VisitFloatType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitFloatType(node *ast.FloatType) error {
	return nil
}

// VisitFunctionCall implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	return nil
}

// VisitFunctionDefinition implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	return nil
}

// VisitMainStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitMainStatement(node *ast.MainStatement) error {
	return nil
}

// VisitMemberExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitMemberExpression(node *ast.MemberExpression) error {
	return nil
}

// VisitNumber64Type implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitNumber64Type(node *ast.Number64Type) error {
	return nil
}

// VisitNumberExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitNumberExpression(node *ast.NumberExpression) error {
	return nil
}

// VisitNumberType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitNumberType(node *ast.NumberType) error {
	return nil
}

// VisitPointerType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitPointerType(node *ast.PointerType) error {
	return nil
}

// VisitPrefixExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitPrefixExpression(node *ast.PrefixExpression) error {
	return nil
}

// VisitPrintStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitPrintStatement(node *ast.PrintStatetement) error {
	return nil
}

// VisitReturnStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitReturnStatement(node *ast.ReturnStatement) error {
	return nil
}

// VisitStringExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitStringExpression(node *ast.StringExpression) error {
	return nil
}

// VisitStringType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitStringType(node *ast.StringType) error {
	return nil
}

// VisitStructDeclaration implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	return nil
}

// VisitStructInitializationExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	return nil
}

// VisitSymbolExpression implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitSymbolExpression(node *ast.SymbolExpression) error {
	return nil
}

// VisitSymbolType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitSymbolType(node *ast.SymbolType) error {
	return nil
}

// VisitVarDeclaration implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	return nil
}

// VisitVoidType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitVoidType(node *ast.VoidType) error {
	return nil
}

// VisitWhileStatement implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) VisitWhileStatement(node *ast.WhileStatement) error {
	return nil
}

// ZeroOfArrayType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfArrayType(node *ast.ArrayType) error {
	return nil
}

// ZeroOfFloatType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfFloatType(node *ast.FloatType) error {
	return nil
}

// ZeroOfNumber64Type implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfNumber64Type(node *ast.Number64Type) error {
	return nil
}

// ZeroOfNumberType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfNumberType(node *ast.NumberType) error {
	return nil
}

// ZeroOfPointerType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfPointerType(node *ast.PointerType) error {
	return nil
}

// ZeroOfStringType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfStringType(node *ast.StringType) error {
	return nil
}

// ZeroOfSymbolType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfSymbolType(node *ast.SymbolType) error {
	return nil
}

// ZeroOfVoidType implements [ast.CodeGenerator].
func (l *LLVMTypeChecker) ZeroOfVoidType(node *ast.VoidType) error {
	return nil
}

var _ ast.CodeGenerator = (*LLVMTypeChecker)(nil)
