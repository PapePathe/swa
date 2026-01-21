package ast

import (
	"swahili/lang/lexer"
)

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
	Tokens    []lexer.Token
}

var _ Expression = (*CallExpression)(nil)

func (expr CallExpression) Accept(g CodeGenerator) error {
	return g.VisitCallExpression(&expr)
}

func (expr CallExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
