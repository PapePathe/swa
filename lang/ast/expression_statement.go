package ast

import (
	"swahili/lang/lexer"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
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

type ZeroExpression struct {
	T      Type
	Tokens []lexer.Token
}

// Accept implements [Expression].
func (z *ZeroExpression) Accept(g CodeGenerator) error {
	return g.VisitZeroExpression(z)
}

// TokenStream implements [Expression].
func (z *ZeroExpression) TokenStream() []lexer.Token {
	return z.Tokens
}

// VisitedSwaType implements [Expression].
func (z *ZeroExpression) VisitedSwaType() Type {
	panic("unimplemented")
}

var _ Expression = (*ZeroExpression)(nil)
