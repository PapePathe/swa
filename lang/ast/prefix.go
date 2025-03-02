package ast

import "swahili/lang/lexer"

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (n PrefixExpression) expression() {}
