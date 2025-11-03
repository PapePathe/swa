package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type FunctionCallExpression struct {
	Name   Expression
	Args   []Expression
	Tokens []lexer.Token
}

var _ Expression = (*FunctionCallExpression)(nil)

func (expr FunctionCallExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	name, ok := expr.Name.(SymbolExpression)

	if !ok {
		return fmt.Errorf("Expression %v is not a symbol expression", expr.Name), nil
	}

	funcDef := ctx.Module.NamedFunction(name.Value)

	if funcDef.IsNil() {
		return fmt.Errorf("function %s does not exist", name.Value), nil
	}

	err, funcType := ctx.FindFuncSymbol(name.Value)
	if err != nil {
		return fmt.Errorf("functype not defined"), nil
	}

	args := []llvm.Value{}
	for _, arg := range expr.Args {
		err, argVal := arg.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
		args = append(args, *argVal.Value)
	}

	returnValue := ctx.Builder.CreateCall(
		*funcType,
		funcDef,
		args,
		"",
	)
	return nil, &CompilerResult{Value: &returnValue}
}

func (expr FunctionCallExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (fd FunctionCallExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Name"] = fd.Name
	m["Args"] = fd.Args
	m["Tokens"] = fd.Tokens

	res := make(map[string]any)
	res["ast.FunctionCallExpression"] = m

	return json.Marshal(res)
}
