package ast

import "swahili/lang/values"

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

var _ Expression = (*SymbolExpression)(nil)

func (n SymbolExpression) expression() {}

func (v SymbolExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
