package ast

import "swahili/lang/values"

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

var _ Expression = (*CallExpression)(nil)

func (n CallExpression) expression() {}

func (v CallExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
