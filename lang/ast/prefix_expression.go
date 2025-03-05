package ast

import (
	"swahili/lang/lexer"
	"swahili/lang/values"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (n PrefixExpression) expression() {}

func (v PrefixExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
