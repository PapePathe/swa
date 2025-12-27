package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type PrintStatetement struct {
	Values []Expression
	Tokens []lexer.Token
}

var _ Statement = (*PrintStatetement)(nil)

func (ps PrintStatetement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	printableValues := []llvm.Value{}

	for _, v := range ps.Values {
		err, res := v.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch v.(type) {
		case MemberExpression:
			loadedval := ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, "")
			printableValues = append(printableValues, loadedval)
		case ArrayAccessExpression:
			vv, _ := v.(ArrayAccessExpression)

			err, res := vv.CompileLLVMForPrint(ctx)
			if err != nil {
				return err, nil
			}

			printableValues = append(printableValues, *res)
		case ArrayOfStructsAccessExpression:
			printableValues = append(printableValues, *res.Value)
		case NumberExpression, FloatExpression:
			printableValues = append(printableValues, *res.Value)
		case SymbolExpression:
			printableValues = append(printableValues, *res.Value)
		case FunctionCallExpression, BinaryExpression:
			printableValues = append(printableValues, *res.Value)
		case StringExpression:
			glob := llvm.AddGlobal(*ctx.Module, res.Value.Type(), "")
			glob.SetInitializer(*res.Value)

			printableValues = append(printableValues, glob)
		default:
			err := fmt.Errorf("Expression %v not supported in print statement", v)

			return err, nil
		}
	}

	ctx.Builder.CreateCall(
		llvm.FunctionType(llvm.GlobalContext().Int32Type(), []llvm.Type{llvm.PointerType(ctx.Context.Int8Type(), 0)}, true),
		ctx.Module.NamedFunction("printf"),
		printableValues,
		"",
	)

	return nil, nil
}

func (ps PrintStatetement) Accept(g CodeGenerator) error {
	return g.VisitPrintStatement(&ps)
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
