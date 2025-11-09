package ast

import (
	"encoding/json"
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

func (expr AssignmentExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, val := expr.Value.CompileLLVM(ctx)

	if err != nil {
		return err, nil
	}

	err, assignee := expr.Assignee.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if assignee.SymbolTableEntry != nil && assignee.SymbolTableEntry.Address != nil {
		str := ctx.Builder.CreateStore(*val.Value, *assignee.SymbolTableEntry.Address)

		return nil, &CompilerResult{Value: &str}
	}

	panic("test")

	// str := ctx.Builder.CreateStore(*val.Value, *assignee.Value)
	//
	// return nil, &CompilerResult{Value: &str}
}

func (expr AssignmentExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr AssignmentExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Operator"] = expr.Operator
	m["Assignee"] = expr.Assignee
	m["Value"] = expr.Value

	res := make(map[string]any)
	res["ast.AssignmentExpression"] = m

	return json.Marshal(res)
}
