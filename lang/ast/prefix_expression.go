package ast

import (
	"encoding/json"
	"fmt"

	"swahili/lang/lexer"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
	Tokens          []lexer.Token
}

var _ Expression = (*PrefixExpression)(nil)

func (PrefixExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	panic("PrefixExpression compilation is not implemented")
}

func (e PrefixExpression) String() string {
	return fmt.Sprintf("%s %s", e.Operator.Value, e.RightExpression)
}

func (expr PrefixExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs PrefixExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Operator"] = cs.Operator
	m["RightExpression"] = cs.RightExpression
	m["Tokens"] = cs.Tokens

	res := make(map[string]any)
	res["ast.PrefixExpression"] = m

	return json.Marshal(res)
}
