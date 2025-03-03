package ast

import (
	"swahili/lang/lexer"
	"swahili/lang/values"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (n BinaryExpression) expression() {}

func (v BinaryExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
