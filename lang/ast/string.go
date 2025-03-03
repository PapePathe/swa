package ast

import "swahili/lang/values"

// StringExpression ...
type StringExpression struct {
	Value string
}

var _ Expression = (*StringExpression)(nil)

func (n StringExpression) expression() {}

func (v StringExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
