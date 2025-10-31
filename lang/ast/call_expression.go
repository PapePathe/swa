package ast

import "encoding/json"

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

var _ Expression = (*CallExpression)(nil)

func (CallExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	panic("CallExpression compilation is not implemented")
}

func (expr CallExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Caller"] = expr.Caller
	m["Arguments"] = expr.Arguments

	res := make(map[string]any)
	res["ast.CallExpression"] = m

	return json.Marshal(res)
}
