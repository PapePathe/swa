package ast

import "tinygo.org/x/go-llvm"

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

var _ Expression = (*CallExpression)(nil)

func (CallExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
