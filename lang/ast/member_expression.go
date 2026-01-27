package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
	Tokens   []lexer.Token
	SwaType  Type
}

var _ Expression = (*MemberExpression)(nil)

func (expr MemberExpression) Accept(g CodeGenerator) error {
	return g.VisitMemberExpression(&expr)
}

func (expr MemberExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr MemberExpression) String() string {
	return fmt.Sprintf("%s.%s", expr.Object, expr.Property)
}

func (expr MemberExpression) VisitedSwaType() Type {
	return expr.SwaType
}
