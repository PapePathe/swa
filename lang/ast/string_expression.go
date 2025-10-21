package ast

import (
	"encoding/json"

	"tinygo.org/x/go-llvm"
)

// StringExpression ...
type StringExpression struct {
	Value string
}

var _ Expression = (*StringExpression)(nil)

func (se StringExpression) String() string {
	return se.Value
}

func (se StringExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := ctx.Context.ConstString(se.Value, true)

	return nil, &res
}

func (se StringExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value

	res := make(map[string]any)
	res["ast.StringExpression"] = m

	return json.Marshal(res)
}
