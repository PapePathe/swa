package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// IntegerExpression represents an integer literal.
type IntegerExpression struct {
	Value  int64
	Tokens []lexer.Token
}

var _ Expression = (*IntegerExpression)(nil)

func (e IntegerExpression) String() string {
	return fmt.Sprintf("%d", e.Value)
}

func (se IntegerExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), true)

	return nil, &CompilerResult{Value: &res}
}

func (expr IntegerExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (se IntegerExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value

	res := make(map[string]any)
	res["ast.IntegerExpression"] = m

	return json.Marshal(res)
}
