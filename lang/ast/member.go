package ast

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

var _ Expression = (*MemberExpression)(nil)

func (me MemberExpression) expression() {}
