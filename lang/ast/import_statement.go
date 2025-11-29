package ast

import (
	"encoding/json"
	"swahili/lang/lexer"
)

type ImportsStatement struct {
	Files  map[string]any
	Tokens []lexer.Token
}

var _ Statement = (*ImportsStatement)(nil)

func (ps ImportsStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	return nil, nil
}

func (expr ImportsStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr ImportsStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Files"] = expr.Files

	res := make(map[string]any)
	res["ast.ImportsStatetement"] = m

	return json.Marshal(res)
}
