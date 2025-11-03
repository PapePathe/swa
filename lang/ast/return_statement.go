package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type ReturnStatement struct {
	Value  Expression
	Tokens []lexer.Token
}

func (rs ReturnStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	switch v := rs.Value.(type) {
	case ArrayAccessExpression:
		err, ptr := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		val := ctx.Builder.CreateLoad(ctx.Context.Int32Type(), *ptr.Value, "")
		ctx.Builder.CreateRet(val)
	case MemberExpression:
		expr, _ := rs.Value.(MemberExpression)

		err, loadedval := expr.CompileLLVMForPropertyAccess(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*loadedval)
	case SymbolExpression:
		err, val := ctx.FindSymbol(v.Value)
		if err != nil {
			return fmt.Errorf("Undefined variable %s", v.Value), nil
		}

		switch val.Value.Type() {
		case llvm.GlobalContext().Int32Type():
			ctx.Builder.CreateRet(val.Value)
		case llvm.PointerType(llvm.GlobalContext().Int32Type(), 0):
			loadedval := ctx.Builder.CreateLoad(llvm.GlobalContext().Int32Type(), val.Value, "")
			ctx.Builder.CreateRet(loadedval)
		}
	case NumberExpression:
		ctx.Builder.CreateRet(llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(v.Value), false))
	case BinaryExpression:
		err, res := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*res.Value)
	default:
		err := fmt.Errorf("unknown expression in ReturnStatement <%s>", rs.Value)

		panic(err)
	}

	return nil, nil
}

func (expr ReturnStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (rs ReturnStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = rs.Value
	m["Tokens"] = rs.Tokens

	res := make(map[string]any)
	res["ast.ReturnStatement"] = m

	return json.Marshal(res)
}
