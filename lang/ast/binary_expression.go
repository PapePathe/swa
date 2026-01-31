package ast

import (
	"fmt"
	"swahili/lang/lexer"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
	Tokens   []lexer.Token
	SwaType  Type
}

var _ Expression = (*BinaryExpression)(nil)

func (expr *BinaryExpression) Accept(g CodeGenerator) error {
	return g.VisitBinaryExpression(expr)
}

func (expr BinaryExpression) String() string {
	return fmt.Sprintf("%s%s%s", expr.Left, expr.Operator.Value, expr.Right)
}

func (expr BinaryExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr BinaryExpression) VisitedSwaType() Type {
	return expr.SwaType
}
