package ast

import (
	"encoding/json"
	"swahili/lang/lexer"
	"swahili/lang/values"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (n PrefixExpression) expression() {}

func (v PrefixExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (cs PrefixExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["operator"] = cs.Operator
	m["right_expression"] = cs.RightExpression

	res := make(map[string]any)
	res["ast.PrefixExpression"] = m

	return json.Marshal(res)
}
