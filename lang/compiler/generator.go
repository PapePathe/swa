package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type LLVMGenerator struct {
	Ctx        *ast.CompilerCtx
	lastResult *ast.CompilerResult
}

// VisitArrayInitializationExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	panic("unimplemented")
}

// VisitArrayOfStructsAccessExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	panic("unimplemented")
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	panic("unimplemented")
}

// VisitBinaryExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitBinaryExpression(node *ast.BinaryExpression) error {
	panic("unimplemented")
}

// VisitBlockStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, v := range node.Body {
		v.Accept(g)
	}

	return nil
}

// VisitCallExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitCallExpression(node *ast.CallExpression) error {
	panic("unimplemented")
}

// VisitConditionalStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	panic("unimplemented")
}

// VisitExpressionStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	panic("unimplemented")
}

// VisitFloatExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFloatExpression(node *ast.FloatExpression) error {
	panic("unimplemented")
}

// VisitFunctionCall implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	panic("unimplemented")
}

// VisitFunctionDefinition implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	panic("unimplemented")
}

// VisitMainStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMainStatement(node *ast.MainStatement) error {
	fnType := llvm.FunctionType(
		g.Ctx.Context.Int32Type(),
		[]llvm.Type{},
		false,
	)
	fn := llvm.AddFunction(*g.Ctx.Module, "main", fnType)
	block := g.Ctx.Context.AddBasicBlock(fn, "entry")
	g.Ctx.Builder.SetInsertPointAtEnd(block)

	return node.Body.Accept(g)
}

// VisitMemberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMemberExpression(node *ast.MemberExpression) error {
	panic("unimplemented")
}

// VisitNumberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitNumberExpression(node *ast.NumberExpression) error {
	panic("unimplemented")
}

// VisitPrefixExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrefixExpression(node *ast.PrefixExpression) error {
	panic("unimplemented")
}

// VisitPrintStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	panic("unimplemented")
}

// VisitReturnStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	switch node.Value.(type) {
	case ast.NumberExpression:
		ret := llvm.ConstInt(g.Ctx.Context.Int32Type(), uint64(node.Value.(ast.NumberExpression).Value), false)
		g.Ctx.Builder.CreateRet(ret)
	default:
		panic("unimplemented")
	}

	return nil
}

// VisitStringExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStringExpression(node *ast.StringExpression) error {
	panic("unimplemented")
}

// VisitStructDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	panic("unimplemented")
}

// VisitStructInitializationExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	panic("unimplemented")
}

// VisitSymbolExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitSymbolExpression(node *ast.SymbolExpression) error {
	panic("unimplemented")
}

// VisitVarDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	panic("unimplemented")
}

// VisitWhileStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitWhileStatement(node *ast.WhileStatement) error {
	panic("unimplemented")
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func NewLLVMGenerator(ctx *ast.CompilerCtx) *LLVMGenerator {
	return &LLVMGenerator{Ctx: ctx}
}

func (g *LLVMGenerator) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	return nil
}
