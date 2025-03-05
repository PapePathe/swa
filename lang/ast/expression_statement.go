package ast

import "swahili/lang/values"

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
