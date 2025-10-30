package ast

import (
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
}

var _ Expression = (*AssignmentExpression)(nil)

func (expr AssignmentExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, val := expr.Value.CompileLLVM(ctx)

	if err != nil {
		return err, nil
	}

	err, assignee := expr.Assignee.CompileLLVM(ctx)

	if err != nil {
		return err, nil
	}

	str := ctx.Builder.CreateStore(*val.Value, *assignee.Value)

	return nil, &CompilerResult{Value: &str}
}
