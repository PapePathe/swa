package ast

import (
	"encoding/json"
	"fmt"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

var _ Expression = (*SymbolExpression)(nil)

func (e SymbolExpression) String() string {
	return e.Value
}

func (se SymbolExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	val, ok := ctx.SymbolTable[se.Value]

	if !ok {
		return fmt.Errorf("Variable %s does not exist", se.Value), nil
	}

	return nil, &CompilerResult{Value: &val.Value}
}

func (expr SymbolExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["operator"] = expr.Value

	res := make(map[string]any)
	res["ast.SymbolExpression"] = m

	return json.Marshal(res)
}
