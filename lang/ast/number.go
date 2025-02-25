package ast

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

func (n NumberExpression) expression() {}
