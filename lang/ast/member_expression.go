package ast

import "tinygo.org/x/go-llvm"

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

var _ Expression = (*MemberExpression)(nil)

func (bs MemberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
