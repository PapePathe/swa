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

	if val.Address != nil {
		// the value may have changed so we load it.
		load := ctx.Builder.CreateLoad(val.Value.Type(), *val.Address, "")

		return nil, &CompilerResult{Value: &load, SymbolTableEntry: val}
	}

	return nil, &CompilerResult{Value: &val.Value, SymbolTableEntry: val}
}

func (expr SymbolExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr SymbolExpression) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.SymbolExpression"] = expr.Value

	return json.Marshal(res)
}
