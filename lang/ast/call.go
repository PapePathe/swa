package ast

// CallExpression ...
type CallExpression struct {
	Arguments []Expression
	Caller    Expression
}

var _ Expression = (*CallExpression)(nil)

func (n CallExpression) expression() {}
