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

		switch typ := arg.(type) {
		case SymbolExpression:
			switch funcDef.Params()[i].Type().TypeKind() {
			case llvm.IntegerTypeKind, llvm.FloatTypeKind, llvm.DoubleTypeKind:
				args = append(args, *argVal.Value)
			case llvm.PointerTypeKind:
				if argVal.SymbolTableEntry.Ref != nil {
					load := ctx.Builder.CreateLoad(argVal.Value.AllocatedType(), *argVal.Value, "")
					alloca := ctx.Builder.CreateAlloca(argVal.Value.AllocatedType(), "")
					ctx.Builder.CreateStore(load, alloca)
					args = append(args, alloca)

					break
				}

				if argVal.SymbolTableEntry != nil && argVal.SymbolTableEntry.Address != nil {
					args = append(args, *argVal.SymbolTableEntry.Address)
				} else {
					args = append(args, *argVal.Value)
				}
			default:
				return fmt.Errorf("Type of%s not supported a func call arg", typ), nil
			}
		default:
			args = append(args, *argVal.Value)
		}

		if i >= len(args) {
			return fmt.Errorf("Arg at position %d does not exist for function %s (%v)", i, name.Value, expr.TokenStream()), nil
		}

		currentArgType := args[i].Type()
		currentParamType := funcDef.Params()[i].Type()

		if currentArgType != currentParamType {
			format := "expected argument of type %v expected but got %v"

			return fmt.Errorf(format, currentParamType, currentArgType), nil
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

func (expr FunctionCallExpression) Accept(g CodeGenerator) error {
	return g.VisitFunctionCall(&expr)
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
