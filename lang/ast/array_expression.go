package ast

import (
	"swahili/lang/values"
)

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (l ArrayInitializationExpression) expression() {}

func (v ArrayInitializationExpression) Evaluate(s *Scope) (error, values.Value) {
	contents := []values.Value{}

	for _, elem := range v.Contents {
		_, exprEval := elem.Evaluate(s)
		contents = append(contents, exprEval)
	}

	array := values.ArrayValue{Values: contents}

	return nil, array
}
