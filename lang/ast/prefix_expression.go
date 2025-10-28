package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (PrefixExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	panic("PrefixExpression compilation is not implemented")
}

func (e PrefixExpression) String() string {
	return fmt.Sprintf("%s %s", e.Operator.Value, e.RightExpression)
}

func (cs PrefixExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["operator"] = cs.Operator
	m["right_expression"] = cs.RightExpression

	res := make(map[string]any)
	res["ast.PrefixExpression"] = m

	return json.Marshal(res)
}
