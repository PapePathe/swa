package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
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
		var load llvm.Value

		switch val.DeclaredType.(type) {
		case StringType:
			load = ctx.Builder.CreateLoad(val.Address.Type(), *val.Address, "")
		case PointerType:
			load = ctx.Builder.CreateLoad(val.Address.AllocatedType(), *val.Address, "")
		case ArrayType:
			load = ctx.Builder.CreateLoad(val.Address.Type(), *val.Address, "")
		case NumberType, FloatType:
			load = ctx.Builder.CreateLoad(val.Address.AllocatedType(), *val.Address, "")
		default:
			load = ctx.Builder.CreateLoad(val.Address.AllocatedType(), *val.Address, "")
		}

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
