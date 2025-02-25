package ast

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

func (l ArrayInitializationExpression) expression() {}
