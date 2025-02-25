package ast

// StringExpression ...
type StringExpression struct {
	Value string
}

func (n StringExpression) expression() {}
