package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

var _ Expression = (*NumberExpression)(nil)

func (e NumberExpression) String() string {
	return fmt.Sprintf("%d", int(e.Value))
}

func (se NumberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), false)

	return nil, &res
}
