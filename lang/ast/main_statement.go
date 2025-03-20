package ast

import (
	"encoding/json"
	"swahili/lang/values"
)

type MainStatement struct {
	Body BlockStatement
}

func (ms MainStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (MainStatement) statement() {}

func (cs MainStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["body"] = cs.Body

	res := make(map[string]any)
	res["ast.MainProGram"] = m

	return json.Marshal(res)
}
