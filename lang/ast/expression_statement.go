package ast

import (
	"encoding/json"

	"swahili/lang/lexer"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp    Expression
	Tokens []lexer.Token
}

var _ Statement = (*ExpressionStatement)(nil)

func (exp ExpressionStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, _ := exp.Exp.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	return nil, nil
}

func (expr ExpressionStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (es ExpressionStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Expression"] = es.Exp

	res := make(map[string]any)
	res["ast.ExpressionStatement"] = m

	return json.Marshal(res)
}
