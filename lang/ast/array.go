package ast

import "swahili/lang/values"

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (l ArrayInitializationExpression) expression() {}

func (v ArrayInitializationExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
