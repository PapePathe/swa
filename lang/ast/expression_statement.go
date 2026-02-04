package ast

import (
	"swahili/lang/lexer"
)

type ExpressionStatement struct {
	Exp    Expression
	Tokens []lexer.Token
}

var _ Statement = (*ExpressionStatement)(nil)

func (expr *ExpressionStatement) Accept(g CodeGenerator) error {
	return g.VisitExpressionStatement(expr)
}

func (expr ExpressionStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
