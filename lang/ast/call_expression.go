package ast

import (
	"encoding/json"

	"swahili/lang/lexer"
)

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
	Tokens    []lexer.Token
}

var _ Expression = (*CallExpression)(nil)

func (CallExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	panic("CallExpression compilation is not implemented")
}

func (expr CallExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr CallExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Caller"] = expr.Caller
	m["Arguments"] = expr.Arguments

	res := make(map[string]any)
	res["ast.CallExpression"] = m

	return json.Marshal(res)
}
