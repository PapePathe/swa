package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitErrorExpression(node *ast.ErrorExpression) error {
	// Treat ‘error’ keyword as value 1 (error present)
	val := llvm.ConstInt(g.Ctx.Context.Int32Type(), 1, false)
	res := &CompilerResult{
		Value: &val,
	}
	g.setLastResult(res)
	return nil
}

func (g *LLVMGenerator) VisitErrorType(node *ast.ErrorType) error {
	// In LLVM, error could be an i32 (0 for no error, non-zero for error code)
	// or a pointer to a struct. Let's use i32 for now as a simple implementation.
	res := CompilerResultType{
		Type: g.Ctx.Context.Int32Type(),
	}
	g.setLastTypeVisitResult(&res)
	return nil
}

func (g *LLVMGenerator) ZeroOfErrorType(node *ast.ErrorType) error {
	// Zero of error type is 'no error' (0)
	val := llvm.ConstInt(g.Ctx.Context.Int32Type(), 0, false)
	g.setLastResult(&CompilerResult{
		Value: &val,
	})
	return nil
}
