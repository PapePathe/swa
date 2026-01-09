package ast

import (
	"encoding/json"

	"swahili/lang/lexer"
)

// StringExpression ...
type StringExpression struct {
	Value  string
	Tokens []lexer.Token
}

var _ Expression = (*StringExpression)(nil)

func (expr StringExpression) String() string {
	return expr.Value
}

func (expr StringExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	res := ctx.Context.ConstString(expr.Value, true)

	return nil, &CompilerResult{Value: &res}
}

func (expr StringExpression) Accept(g CodeGenerator) error {
	return g.VisitStringExpression(&expr)
}

func (expr StringExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr StringExpression) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.StringExpression"] = expr.Value

	return json.Marshal(res)
}
