package ast

import (
	"swahili/lang/lexer"
)

type ArrayAccessExpression struct {
	Name   Expression
	Index  Expression
	Tokens []lexer.Token
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr ArrayAccessExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayAccessExpression(&expr)
}

func (expr ArrayAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
