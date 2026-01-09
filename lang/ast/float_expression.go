package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// FloatExpression represents a floating-point literal.
type FloatExpression struct {
	Value  float64
	Tokens []lexer.Token
}

var _ Expression = (*FloatExpression)(nil)

func (e FloatExpression) String() string {
	return fmt.Sprintf("%f", e.Value)
}

func (se FloatExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	res := llvm.ConstFloat(llvm.GlobalContext().DoubleType(), se.Value)

	return nil, &CompilerResult{Value: &res}
}

func (expr FloatExpression) Accept(g CodeGenerator) error {
	return g.VisitFloatExpression(&expr)
}

func (expr FloatExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (se FloatExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value

	res := make(map[string]any)
	res["ast.FloatExpression"] = m

	return json.Marshal(res)
}
