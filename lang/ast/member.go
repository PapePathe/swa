package ast

import "swahili/lang/values"

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

var _ Expression = (*MemberExpression)(nil)

func (me MemberExpression) expression() {}

func (v MemberExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
