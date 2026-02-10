package compiler

import (
	"swahili/lang/ast"
)

func (g *LLVMGenerator) VisitErrorExpression(node *ast.ErrorExpression) error {
	err := node.Exp.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastResult()

	g.Debugf("VisitErrorExpression %v", lastres)

	res := &CompilerResult{
		Value: lastres.Value,
		SymbolTableEntry: &SymbolTableEntry{
			Address:      lastres.Value,
			DeclaredType: ast.ErrorType{},
		},
		SwaType: ast.ErrorType{},
	}

	g.setLastResult(res)

	return nil
}
