package ast

import (
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
