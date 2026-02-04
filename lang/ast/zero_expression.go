package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

type ZeroExpression struct {
	T      Type
	Tokens []lexer.Token
}

var _ Expression = (*ZeroExpression)(nil)

func (z *ZeroExpression) Accept(g CodeGenerator) error {
	return g.VisitZeroExpression(z)
}

func (z *ZeroExpression) TokenStream() []lexer.Token {
	return z.Tokens
}

func (z *ZeroExpression) VisitedSwaType() Type {
	return z.T
}

func (z ZeroExpression) String() string {
	return fmt.Sprintf("zero of %s", z.T.Value().String())
}
