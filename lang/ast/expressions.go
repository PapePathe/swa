package ast

import "swahili/lang/lexer"

type NumberExpression struct {
	Value float64
}

func (n NumberExpression) expression() {

}

type StringExpression struct {
	Value string
}

func (n StringExpression) expression() {

}

type SymbolExpression struct {
	Value string
}

func (n SymbolExpression) expression() {

}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

func (n BinaryExpression) expression() {

}
