package ast

import (
	"fmt"

	"swahili/lang/lexer"
)

// NumberExpression ...
type NumberExpression struct {
	Value   int64
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*NumberExpression)(nil)

func (e NumberExpression) String() string {
	return fmt.Sprintf("%d", e.Value)
}

func (expr *NumberExpression) Accept(g CodeGenerator) error {
	return g.VisitNumberExpression(expr)
}

func (expr NumberExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr NumberExpression) VisitedSwaType() Type {
	return NumberType{}
}
