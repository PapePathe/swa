package ast

type StructInitializationExpression struct {
	Name       string
	Properties map[string]Expression
}

var _ Expression = (*StructInitializationExpression)(nil)

func (n StructInitializationExpression) expression() {}
