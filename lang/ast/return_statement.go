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

		val := ctx.Builder.CreateLoad(llvm.GlobalContext().Int32Type(), *ptr.Value, "")
		ctx.Builder.CreateRet(val)
	case MemberExpression:
		expr, _ := rs.Value.(MemberExpression)

		err, loadedval := expr.CompileLLVMForPropertyAccess(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*loadedval)
	case SymbolExpression:
		err, val := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*val.Value)
	case StringExpression:
		err, ptr := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		alloc := ctx.Builder.CreateAlloca(ptr.Value.Type(), "")
		ctx.Builder.CreateStore(*ptr.Value, alloc)
		ctx.Builder.CreateRet(alloc)

	case NumberExpression:
		ctx.Builder.CreateRet(llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(v.Value), false))
	case FloatExpression:
		val := rs.Value.(FloatExpression)
		ctx.Builder.CreateRet(llvm.ConstFloat(llvm.GlobalContext().DoubleType(), val.Value))
	case BinaryExpression:
		err, res := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*res.Value)
	case FunctionCallExpression:
		err, res := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		ctx.Builder.CreateRet(*res.Value)
	default:
		err := fmt.Errorf("unknown expression in ReturnStatement <%s>", rs.Value)
		return err, nil
	}

	return nil, nil
}

func (expr ReturnStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (rs ReturnStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = rs.Value

	res := make(map[string]any)
	res["ast.ReturnStatement"] = m

	return json.Marshal(res)
}
