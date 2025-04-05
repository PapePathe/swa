package ast

import (
	"encoding/json"

	"github.com/llir/llvm/ir/types"
	"tinygo.org/x/go-llvm"
)

type MainStatement struct {
	Body BlockStatement
}

func (ms MainStatement) Compile(ctx *Context) error {
	main := ctx.mod.NewFunc("main", types.I32)
	mainCtx := ctx.NewContext(main.NewBlock(""))

	return ms.Body.Compile(mainCtx)
}

func (ms MainStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
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

func (cs MainStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["body"] = cs.Body

	res := make(map[string]any)
	res["ast.MainProGram"] = m

	return json.Marshal(res)
}
