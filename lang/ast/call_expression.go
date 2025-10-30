package ast

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

var _ Expression = (*CallExpression)(nil)

func (CallExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	panic("CallExpression compilation is not implemented")
}
