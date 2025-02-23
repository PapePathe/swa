package ast

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

func (me MemberExpression) expression() {}
