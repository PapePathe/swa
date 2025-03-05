package ast

import (
	"swahili/lang/lexer"
	"swahili/lang/values"
)

// AssignmentExpression.
// Is an expression where the programmer is trying to assign a value to a variable.
//
// a = a +5;
// foo.bar = foo.bar + 10;
type AssignmentExpression struct {
	Operator lexer.Token
	Assignee Expression
	Value    Expression
}

var _ Expression = (*AssignmentExpression)(nil)

func (n AssignmentExpression) expression() {}

func (v AssignmentExpression) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
