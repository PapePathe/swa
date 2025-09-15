package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

var _ Expression = (*SymbolExpression)(nil)

func (se SymbolExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	val, ok := ctx.SymbolTable[se.Value]

	if !ok {
		return fmt.Errorf("Variable %s does not exist", se.Value), nil
	}

	return nil, &val.Value
}
