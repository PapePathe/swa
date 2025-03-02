package ast

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (l ArrayInitializationExpression) expression() {}
