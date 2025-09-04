package ast

import (
	"tinygo.org/x/go-llvm"
)

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

var _ Expression = (*NumberExpression)(nil)

func (se NumberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), false)

	return nil, &res
}
