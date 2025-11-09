package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

type PrintStatetement struct {
	Values []Expression
	Tokens []lexer.Token
}

var _ Statement = (*PrintStatetement)(nil)

// TODO : print all the values in a single call
func (ps PrintStatetement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	printableValues := []llvm.Value{}

	for _, v := range ps.Values {
		err, res := v.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch v.(type) {
		case MemberExpression:
			expr, _ := v.(MemberExpression)

			err, loadedval := expr.CompileLLVMForPropertyAccess(ctx)
			if err != nil {
				return err, nil
			}
			printableValues = append(printableValues, *loadedval)
		case ArrayAccessExpression:
			vv, _ := v.(ArrayAccessExpression)
			err, res := vv.CompileLLVMForPrint(ctx)
			if err != nil {
				return err, nil
			}
			printableValues = append(printableValues, *res)
		case ArrayOfStructsAccessExpression:
			printableValues = append(printableValues, *res.Value)
		case NumberExpression:
			printableValues = append(printableValues, *res.Value)
		case SymbolExpression:
			printableValues = append(printableValues, *res.Value)
		case StringExpression:
			all := ctx.Builder.CreateAlloca(res.Value.Type(), "")
			ctx.Builder.CreateStore(*res.Value, all)

			printableValues = append(printableValues, all)
		default:
			panic(fmt.Sprintf("Expression %v not supported in print statement", v))
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

func (expr PrintStatetement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr PrintStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Values"] = expr.Values

	res := make(map[string]any)
	res["ast.PrintStatetement"] = m

	return json.Marshal(res)
}
