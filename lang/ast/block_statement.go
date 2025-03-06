package ast

import (
	"encoding/json"
	"swahili/lang/values"
)

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body []Statement
}

var _ Statement = (*BlockStatement)(nil)

func (cs BlockStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (bs BlockStatement) statement() {}

func (bs BlockStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["ast.BlockStatement"] = bs.Body

	return json.Marshal(m)
}
