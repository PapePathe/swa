package ast

import "swahili/lang/lexer"

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (n BinaryExpression) expression() {}
