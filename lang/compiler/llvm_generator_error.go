package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitErrorExpression(node *ast.ErrorExpression) error {
	err := node.Exp.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastResult()

	g.Debugf("VisitErrorExpression %+v", lastres)

	val := llvm.ConstInt(g.Ctx.Context.Int32Type(), 1, false)
	res := &CompilerResult{
		Value: &val,
		SymbolTableEntry: &SymbolTableEntry{
			Address:      lastres.Value,
			DeclaredType: ast.ErrorType{},
		},
	}

	g.setLastResult(res)

	return nil
}

func (g *LLVMGenerator) VisitErrorType(node *ast.ErrorType) error {
	res := CompilerResultType{
		Type: g.Ctx.Context.Int32Type(),
	}
	g.setLastTypeVisitResult(&res)
	return nil
}

func (g *LLVMGenerator) ZeroOfErrorType(node *ast.ErrorType) error {
	val := llvm.ConstInt(g.Ctx.Context.Int32Type(), 0, false)
	g.setLastResult(&CompilerResult{
		Value: &val,
	})

	return nil
}
