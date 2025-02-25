package ast

type StructInitializationExpression struct {
	Name       string
	Properties map[string]Expression
}

func (n StructInitializationExpression) expression() {}
