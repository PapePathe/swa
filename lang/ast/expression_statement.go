package ast

import (
	"encoding/json"
	"swahili/lang/values"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp Expression
}

var _ Statement = (*ExpressionStatement)(nil)

func (cs ExpressionStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (bs ExpressionStatement) statement() {}

func (cs ExpressionStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["expression"] = cs.Exp

	res := make(map[string]any)
	res["ast.ExpressionStatetement"] = m

	return json.Marshal(res)
}
