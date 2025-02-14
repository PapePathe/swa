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

func (n BinaryExpression) expression() {

}
