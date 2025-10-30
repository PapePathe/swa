package ast

import (
	"encoding/json"
)

// StringExpression ...
type StringExpression struct {
	Value string
}

var _ Expression = (*StringExpression)(nil)

func (se StringExpression) String() string {
	return se.Value
}

func (se StringExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	res := ctx.Context.ConstString(se.Value, true)

	return nil, &CompilerResult{Value: &res}
}

func (se StringExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value

	res := make(map[string]any)
	res["ast.StringExpression"] = m

	return json.Marshal(res)
}
