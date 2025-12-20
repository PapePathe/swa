package ast

import (
	"encoding/json"
	"fmt"

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
	return fmt.Errorf("CallExpression compilation is not implemented"), nil
}

func (expr CallExpression) Accept(g CodeGenerator) error {
	return g.VisitCallExpression(&expr)
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
