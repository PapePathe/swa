package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value  string
	Tokens []lexer.Token
}

var _ Expression = (*SymbolExpression)(nil)

func (e SymbolExpression) String() string {
	return e.Value
}

func (expr SymbolExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, val := ctx.FindSymbol(expr.Value)
	if err != nil {
		return fmt.Errorf("Variable %s does not exist", expr.Value), nil
	}

	return nil, &CompilerResult{Value: &val.Value, SymbolTableEntry: val}
}

func (expr SymbolExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr SymbolExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = expr.Value

	res := make(map[string]any)
	res["ast.SymbolExpression"] = m

	return json.Marshal(res)
}
