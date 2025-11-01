package ast

import (
	"encoding/json"
	"fmt"
	"os"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
	Tokens   []lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (expr BinaryExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, compiledLeftValue, compiledRightValue := expr.compileLeftAndRight(ctx)
	if err != nil {
		return err, nil
	}

	var finalLeftValue, finalRightValue llvm.Value

	ctype := commonType(*compiledLeftValue, *compiledRightValue)

	if compiledLeftValue.Type().TypeKind() == llvm.PointerTypeKind {
		finalLeftValue = ctx.Builder.CreateLoad(compiledLeftValue.GlobalValueType(), *compiledLeftValue, "")
	} else {
		finalLeftValue = expr.castToType(ctx, ctype, *compiledLeftValue)
	}

	if compiledRightValue.Type().TypeKind() == llvm.PointerTypeKind {
		finalRightValue = ctx.Builder.CreateLoad(compiledRightValue.GlobalValueType(), *compiledRightValue, "")
	} else {
		finalRightValue = expr.castToType(ctx, ctype, *compiledRightValue)
	}

	handler, ok := handlers[expr.Operator.Kind]
	if !ok {
		return fmt.Errorf("unsupported operator <%s>", expr.Operator.Kind), nil
	}
	return handler(ctx, finalLeftValue, finalRightValue)
}

func (expr BinaryExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr BinaryExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Left"] = expr.Left
	m["Right"] = expr.Right
	m["Operator"] = expr.Operator
	m["Tokens"] = expr.TokenStream()

	res := make(map[string]any)
	res["ast.BinaryExpression"] = m

	return json.Marshal(res)
}

func (expr BinaryExpression) compileLeftAndRight(ctx *CompilerCtx) (error, *llvm.Value, *llvm.Value) {
	err, compiledLeftValue := expr.Left.CompileLLVM(ctx)
	if err != nil {
		return err, nil, nil
	}

	if compiledLeftValue == nil {
		return fmt.Errorf("left side of expression is nil"), nil, nil
	}

	if compiledLeftValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("left side of expression is of type void"), nil, nil
	}

	err, compiledRightValue := expr.Right.CompileLLVM(ctx)
	if err != nil {
		return err, nil, nil
	}

	if compiledRightValue == nil {
		return fmt.Errorf("right side of expression is nil"), nil, nil
	}

	if compiledRightValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("right side of expression is of type void"), nil, nil
	}
	return nil, compiledLeftValue.Value, compiledRightValue.Value
}

func commonType(l, r llvm.Value) llvm.Type {
	if l.Type().TypeKind() == llvm.PointerTypeKind && l.Type().TypeKind() == llvm.PointerTypeKind {
		return l.GlobalValueType()
	}

	if l.Type() == r.Type() {
		return l.Type()
	}
	if l.Type().TypeKind() == llvm.PointerTypeKind {
		return l.GlobalValueType()
	}

	if r.Type().TypeKind() == llvm.PointerTypeKind {
		return r.GlobalValueType()
	}

	panic(fmt.Errorf("Unhandled combination %v %v", r.Type(), l.Type()))
}

func (expr BinaryExpression) castToType(ctx *CompilerCtx, t llvm.Type, v llvm.Value) llvm.Value {
	if t.TypeKind() == v.Type().TypeKind() {
		return v
	}

	switch v.Type().TypeKind() {
	case llvm.PointerTypeKind:
		switch t.TypeKind() {
		case llvm.IntegerTypeKind:
			return ctx.Builder.CreatePtrToInt(v, t, "")
		}
	case llvm.IntegerTypeKind:
		fmt.Println(fmt.Errorf("Value %v is of type integer and other is %s", v.Type().TypeKind(), t.TypeKind()))
		os.Exit(1)
	default:
		panic(fmt.Errorf("Unhandled type %s, %s", t.TypeKind(), v.Type().TypeKind()))
	}

	return llvm.Value{}
}

type binaryHandlerFunc func(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult)

var handlers = map[lexer.TokenKind]binaryHandlerFunc{
	lexer.Plus:              add,
	lexer.GreaterThan:       greaterThan,
	lexer.GreaterThanEquals: greaterThanEquals,
	lexer.LessThan:          lessThan,
	lexer.LessThanEquals:    lessThanEquals,
	lexer.Equals:            equals,
}

func add(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateAdd(l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func greaterThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateICmp(llvm.IntUGT, l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func greaterThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateICmp(llvm.IntUGE, l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func lessThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateICmp(llvm.IntULT, l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func lessThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateICmp(llvm.IntULE, l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func equals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateICmp(llvm.IntEQ, l, r, "")

	return nil, &CompilerResult{Value: &res}
}
