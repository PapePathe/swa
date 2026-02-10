package ast

import (
	"fmt"

	"swahili/lang/lexer"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
	Tokens          []lexer.Token
	SwaType         Type
}

var _ Expression = (*PrefixExpression)(nil)

func (e PrefixExpression) String() string {
	return fmt.Sprintf("%s %s", e.Operator.Value, e.RightExpression)
}

func (expr *PrefixExpression) Accept(g CodeGenerator) error {
	return g.VisitPrefixExpression(expr)
}

func (expr PrefixExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr PrefixExpression) VisitedSwaType() Type {
	return expr.SwaType
}
