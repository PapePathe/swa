package ast

import "swahili/lang/values"

type ReturnStatement struct {
	Value Expression
}

func (rs ReturnStatement) statement() {}
func (rs ReturnStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
