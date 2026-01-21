package ast

import (
	"fmt"
	"swahili/lang/lexer"
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
	Tokens   []lexer.Token
}

var _ Expression = (*AssignmentExpression)(nil)

func (expr AssignmentExpression) String() string {
	return fmt.Sprintf("%s %s %s", expr.Assignee, expr.Operator.Value, expr.Value)
}

func (expr AssignmentExpression) Accept(g CodeGenerator) error {
	return g.VisitAssignmentExpression(&expr)
}

func (expr AssignmentExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
