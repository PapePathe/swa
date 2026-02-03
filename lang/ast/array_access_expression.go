package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

type ArrayAccessExpression struct {
	Name    Expression
	Index   Expression
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*ArrayAccessExpression)(nil)

func (expr *ArrayAccessExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayAccessExpression(expr)
}

func (expr ArrayAccessExpression) String() string {
	return fmt.Sprintf("%s[%s]", expr.Name, expr.Index)
}
func (expr ArrayAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr ArrayAccessExpression) VisitedSwaType() Type {
	return expr.SwaType
}
