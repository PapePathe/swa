package ast

import (
	"fmt"

	"swahili/lang/lexer"
)

// FloatExpression represents a floating-point literal.
type FloatExpression struct {
	Value   float64
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*FloatExpression)(nil)

func (e FloatExpression) String() string {
	return fmt.Sprintf("%f", e.Value)
}

func (expr *FloatExpression) Accept(g CodeGenerator) error {
	return g.VisitFloatExpression(expr)
}

func (expr FloatExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr FloatExpression) VisitedSwaType() Type {
	return FloatType{}
}
