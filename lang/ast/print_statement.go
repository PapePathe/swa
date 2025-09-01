package ast

import (
	"encoding/json"

	"tinygo.org/x/go-llvm"
)

type PrintStatetement struct {
	Values []Expression
}

var _ Statement = (*PrintStatetement)(nil)

// TODO : print all the values in a single call
func (ps PrintStatetement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	printableValues := []llvm.Value{}

	for _, v := range ps.Values {
		err, res := v.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch v.(type) {
		case SymbolExpression:
			name := v.(SymbolExpression).Value
			global := ctx.Module.NamedGlobal(name)
			//			all := ctx.Builder.CreateLoad(res.Type(), global, "")

			printableValues = append(printableValues, global)

		case StringExpression:
			all := ctx.Builder.CreateAlloca(res.Type(), "")
			ctx.Builder.CreateStore(*res, all)

			printableValues = append(printableValues, all)
		}
	}

	ctx.Builder.CreateCall(
		llvm.FunctionType(ctx.Context.Int32Type(), []llvm.Type{llvm.PointerType(ctx.Context.Int8Type(), 0)}, true),
		ctx.Module.NamedFunction("printf"),
		printableValues,
		"",
	)
	return nil, nil
}

func (cs PrintStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Values"] = cs.Values

	res := make(map[string]any)
	res["ast.PrintStatetement"] = m

	return json.Marshal(res)
}
