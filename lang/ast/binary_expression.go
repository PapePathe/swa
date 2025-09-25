package ast

import (
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (be BinaryExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	err, compiledLeftValue := be.Left.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if compiledLeftValue == nil {
		return fmt.Errorf("left side of expression is nil"), nil
	}

	err, compiledRightValue := be.Right.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if compiledRightValue == nil {
		return fmt.Errorf("right side of expression is nil"), nil
	}

	handler, ok := handlers[be.Operator.Kind]
	if !ok {
		return fmt.Errorf("unsupported operator <%s>", be.Operator.Kind), nil
	}

	var finalLeftValue, finalRightValue llvm.Value

	ctype := commonType(*compiledLeftValue, *compiledRightValue)

	if compiledLeftValue.Type().TypeKind() == llvm.PointerTypeKind {
		finalLeftValue = ctx.Builder.CreateLoad(compiledLeftValue.GlobalValueType(), *compiledLeftValue, "")
	} else {
		finalLeftValue = be.castToType(ctx, ctype, *compiledLeftValue)
	}

	if compiledRightValue.Type().TypeKind() == llvm.PointerTypeKind {
		finalRightValue = ctx.Builder.CreateLoad(compiledRightValue.GlobalValueType(), *compiledRightValue, "")
	} else {
		finalRightValue = be.castToType(ctx, ctype, *compiledRightValue)
	}

	return handler(ctx, finalLeftValue, finalRightValue)
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

	panic(fmt.Errorf("Unhandled combination %v %v", r, l))
}

func (be BinaryExpression) castToType(ctx *CompilerCtx, t llvm.Type, v llvm.Value) llvm.Value {
	if t == v.Type() {
		return v
	}

	switch v.Type().TypeKind() {
	case llvm.PointerTypeKind:
		switch t.TypeKind() {
		case llvm.IntegerTypeKind:
			return ctx.Builder.CreatePtrToInt(v, t, "")
		}
	default:
		panic(fmt.Errorf("Unhandled type %v, %v", v, t))
	}

	return llvm.Value{}
}

type binaryHandlerFunc func(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value)

var handlers = map[lexer.TokenKind]binaryHandlerFunc{
	lexer.Plus:              add,
	lexer.GreaterThan:       greaterThan,
	lexer.GreaterThanEquals: greaterThanEquals,
	lexer.LessThan:          lessThan,
	lexer.LessThanEquals:    lessThanEquals,
	lexer.Equals:            equals,
}

func add(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateAdd(l, r, "")

	return nil, &res
}

func greaterThan(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateICmp(llvm.IntUGT, l, r, "")

	return nil, &res
}

func greaterThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateICmp(llvm.IntUGE, l, r, "")

	return nil, &res
}

func lessThan(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateICmp(llvm.IntULT, l, r, "")

	return nil, &res
}

func lessThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateICmp(llvm.IntULE, l, r, "")

	return nil, &res
}

func equals(ctx *CompilerCtx, l, r llvm.Value) (error, *llvm.Value) {
	res := ctx.Builder.CreateICmp(llvm.IntEQ, l, r, "")

	return nil, &res
}
