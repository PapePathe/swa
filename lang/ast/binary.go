package ast

import "swahili/lang/lexer"

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (n BinaryExpression) expression() {}
