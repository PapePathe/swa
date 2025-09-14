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

func (exp ExpressionStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	err, _ := exp.Exp.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (es ExpressionStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["expression"] = es.Exp

	res := make(map[string]any)
	res["ast.ExpressionStatement"] = m

	return json.Marshal(res)
}
