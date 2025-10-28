package ast

import (
	"fmt"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

var _ Expression = (*SymbolExpression)(nil)

func (e SymbolExpression) String() string {
	return e.Value
}

func (se SymbolExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	val, ok := ctx.SymbolTable[se.Value]

	if !ok {
		return fmt.Errorf("Variable %s does not exist", se.Value), nil
	}

	return nil, &CompilerResult{Value: &val.Value}
}
