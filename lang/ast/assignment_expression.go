package ast

import (
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
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

func (bs AssignmentExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
