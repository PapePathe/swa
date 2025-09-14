package ast

import (
	"encoding/json"

	"tinygo.org/x/go-llvm"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp Expression
}

var _ Statement = (*ExpressionStatement)(nil)

func (ExpressionStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}

func (es ExpressionStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["expression"] = es.Exp

	res := make(map[string]any)
	res["ast.ExpressionStatetement"] = m

	return json.Marshal(res)
}
