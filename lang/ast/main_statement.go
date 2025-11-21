package ast

import (
	"encoding/json"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type MainStatement struct {
	Body   BlockStatement
	Tokens []lexer.Token
}

func (ms MainStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	mainFunc := llvm.AddFunction(
		*ctx.Module,
		"main",
		llvm.FunctionType(
			llvm.GlobalContext().Int32Type(),
			[]llvm.Type{},
			false,
		),
	)

	block := ctx.Context.AddBasicBlock(mainFunc, "func-body")
	ctx.Builder.SetInsertPointAtEnd(block)

	return ms.Body.CompileLLVM(ctx)
}

func (expr MainStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs MainStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Body"] = cs.Body

	res := make(map[string]any)
	res["ast.MainProGram"] = m

	return json.Marshal(res)
}
