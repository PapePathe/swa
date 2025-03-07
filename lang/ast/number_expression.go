package ast

import (
	"encoding/json"
	"swahili/lang/values"
)

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

var _ Expression = (*MemberExpression)(nil)

func (n NumberExpression) expression() {}

func (n NumberExpression) Evaluate(_ *Scope) (error, values.Value) {
	return nil, values.NumberValue{Value: n.Value}
}

func (cs NumberExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["value"] = cs.Value

	res := make(map[string]any)
	res["ast.NumberExpression"] = m

	return json.Marshal(res)
}
