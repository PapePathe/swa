package ast

import (
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

func (ms MainStatement) Accept(g CodeGenerator) error {
	return g.VisitMainStatement(&ms)
}

func (expr MainStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
