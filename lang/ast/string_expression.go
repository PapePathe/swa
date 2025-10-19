package ast

import (
	"tinygo.org/x/go-llvm"
)

// StringExpression ...
type StringExpression struct {
	Value string
}

var _ Expression = (*StringExpression)(nil)

func (se StringExpression) String() string {
	return se.Value
}

func (se StringExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := ctx.Context.ConstString(se.Value, true)

	return nil, &res
}
