package ast

import "swahili/lang/values"

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
func (cs ConditionalStatetement) statement() {}
