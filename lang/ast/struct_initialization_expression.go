package ast

import "tinygo.org/x/go-llvm"

type StructInitializationExpression struct {
	Name       string
	Properties map[string]Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

func (StructInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
