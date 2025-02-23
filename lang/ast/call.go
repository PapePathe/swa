package ast

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

func (n CallExpression) expression() {}
