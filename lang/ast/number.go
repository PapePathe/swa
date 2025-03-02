package ast

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

var _ Expression = (*MemberExpression)(nil)

func (n NumberExpression) expression() {}
