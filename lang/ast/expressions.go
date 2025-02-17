package ast

import "swahili/lang/lexer"

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

func (n NumberExpression) expression() {

}

// StringExpression ...
type StringExpression struct {
	Value string
}

func (n StringExpression) expression() {

}

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

func (n SymbolExpression) expression() {

}

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (n BinaryExpression) expression() {}

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

func (n PrefixExpression) expression() {}

// AssignmentExpression is an expression where the
// programmer is trying to assign a value to a variable.
//
// a = a +5;
// foo.bar = foo.bar + 10;
type AssignmentExpression struct {
	Operator lexer.Token
	Assignee Expression
	Value    Expression
}

func (n AssignmentExpression) expression() {}
