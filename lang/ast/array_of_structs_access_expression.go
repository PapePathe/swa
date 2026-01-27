package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

type ArrayOfStructsAccessExpression struct {
	Name     Expression
	Index    Expression
	Property Expression
	Tokens   []lexer.Token
	SwaType  Type
}

func (expr ArrayOfStructsAccessExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayOfStructsAccessExpression(&expr)
}

func (expr ArrayOfStructsAccessExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr ArrayOfStructsAccessExpression) String() string {
	return fmt.Sprintf("%s[%s].%s", expr.Name, expr.Index, expr.Property)
}

func (expr ArrayOfStructsAccessExpression) VisitedSwaType() Type {
	return expr.SwaType
}
