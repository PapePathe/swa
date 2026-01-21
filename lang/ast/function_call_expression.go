package ast

import (
	"swahili/lang/lexer"
)

type FunctionCallExpression struct {
	Name   Expression
	Args   []Expression
	Tokens []lexer.Token
}

var _ Expression = (*FunctionCallExpression)(nil)

func (expr FunctionCallExpression) Accept(g CodeGenerator) error {
	return g.VisitFunctionCall(&expr)
}

func (expr FunctionCallExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
