package ast

import "swahili/lang/values"

type StructInitializationExpression struct {
	Name       string
	Properties map[string]Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

func (n StructInitializationExpression) expression() {}

func (v StructInitializationExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
