package ast

import (
	"swahili/lang/lexer"
)

type FunctionCallExpression struct {
	Name    Expression
	Args    []Expression
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*FunctionCallExpression)(nil)

func (expr FunctionCallExpression) Accept(g CodeGenerator) error {
	return g.VisitFunctionCall(&expr)
}

func (expr FunctionCallExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr FunctionCallExpression) VisitedSwaType() Type {
	return expr.SwaType
}
