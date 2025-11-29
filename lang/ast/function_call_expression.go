package ast

import (
	"encoding/json"
	"fmt"
	"strings"
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
	var funcName string
	switch expr.Name.(type) {
	case SymbolExpression:
		name, _ := expr.Name.(SymbolExpression)
		funcName = name.Value
	case PackageAccessExpression:
		name, _ := expr.Name.(PackageAccessExpression)
		values := strings.Split(name.Name(), "/")
		funcName = values[1]
	default:
		return fmt.Errorf("Expression %v is not a symbol or package access expression", expr.Name), nil
	}
	funcDef := ctx.Module.NamedFunction(funcName)
	if funcDef.IsNil() {
		return fmt.Errorf("function %s does not exist", funcName), nil
	}

	err, funcType := ctx.FindFuncSymbol(funcName)
	if err != nil {
		return fmt.Errorf("functype not defined"), nil
	}

	argsCount := len(expr.Args)
	paramsCount := len(funcDef.Params())
	if argsCount != paramsCount {
		format := "function %s expect %d arguments but was given %d"
		return fmt.Errorf(format, name.Value, paramsCount, argsCount), nil
	}

	args := []llvm.Value{}
	for i, arg := range expr.Args {
		err, argVal := arg.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch arg.(type) {
		case StringExpression:
			glob := llvm.AddGlobal(*ctx.Module, argVal.Value.Type(), "")
			glob.SetInitializer(*argVal.Value)
			args = append(args, glob)
		default:
			args = append(args, *argVal.Value)
		}
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

	res := make(map[string]any)
	res["ast.FunctionCallExpression"] = m

	return json.Marshal(res)
}
