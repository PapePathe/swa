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
		return rs.compileArrayAccess(ctx)
	case MemberExpression:
		return rs.compileMemberExpression(ctx)
	case SymbolExpression:
		return rs.compileSySymbolExpression(ctx)
	case StringExpression:
		return rs.compileStringExpression(ctx)
	case NumberExpression:
		ctx.Builder.CreateRet(llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(v.Value), false))
	case FloatExpression:
		val, _ := rs.Value.(FloatExpression)
		ctx.Builder.CreateRet(llvm.ConstFloat(llvm.GlobalContext().DoubleType(), val.Value))
	case BinaryExpression:
		return rs.compileBinaryExpression(ctx)
	case FunctionCallExpression:
		return rs.compileFunctionCallExpression(ctx)
	default:
		err := fmt.Errorf("ReturnStatement unknown expression <%s>", rs.Value)
		return err, nil
	}

	return nil, nil
}

func (rs ReturnStatement) Accept(g CodeGenerator) error {
	return g.VisitReturnStatement(&rs)
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

func (rs ReturnStatement) compileArrayAccess(ctx *CompilerCtx) (error, *CompilerResult) {
	err, ptr := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	val := ctx.Builder.CreateLoad(llvm.GlobalContext().Int32Type(), *ptr.Value, "")
	ctx.Builder.CreateRet(val)

	return nil, nil
}

func (rs ReturnStatement) compileMemberExpression(ctx *CompilerCtx) (error, *CompilerResult) {
	err, res := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	loadedval := ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, "")
	ctx.Builder.CreateRet(loadedval)

	return nil, nil
}

func (rs ReturnStatement) compileSySymbolExpression(ctx *CompilerCtx) (error, *CompilerResult) {
	err, val := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateRet(*val.Value)
	return nil, nil
}
func (rs ReturnStatement) compileStringExpression(ctx *CompilerCtx) (error, *CompilerResult) {
	err, ptr := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	alloc := ctx.Builder.CreateAlloca(ptr.Value.Type(), "")
	ctx.Builder.CreateStore(*ptr.Value, alloc)
	ctx.Builder.CreateRet(alloc)

	return nil, nil
}
func (rs ReturnStatement) compileBinaryExpression(ctx *CompilerCtx) (error, *CompilerResult) {
	err, res := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateRet(*res.Value)
	return nil, nil
}
func (rs ReturnStatement) compileFunctionCallExpression(ctx *CompilerCtx) (error, *CompilerResult) {
	err, res := rs.Value.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateRet(*res.Value)
	return nil, nil
}
