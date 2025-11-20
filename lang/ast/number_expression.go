package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// NumberExpression ...
type NumberExpression struct {
	Value  int64
	Tokens []lexer.Token
}

var _ Expression = (*NumberExpression)(nil)

func (e NumberExpression) String() string {
	return fmt.Sprintf("%d", int(e.Value))
}

func (se NumberExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	var signed bool

	if se.Value < 0 {
		signed = true
	}

	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), signed)

	return nil, &CompilerResult{Value: &res}
}

func (expr NumberExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (se NumberExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value

	res := make(map[string]any)
	res["ast.NumberExpression"] = m

	return json.Marshal(res)
}
