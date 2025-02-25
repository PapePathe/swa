package ast

import "swahili/lang/lexer"

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

func (n PrefixExpression) expression() {}
