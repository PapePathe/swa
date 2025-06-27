/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package ast

import (
	"fmt"
	"swahili/lang/lexer"
	"swahili/lang/log"

	"tinygo.org/x/go-llvm"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
}

var (
	_  Expression = (*BinaryExpression)(nil)
	lg            = log.Logger.WithGroup("Ast Evaluator")
)

func (be BinaryExpression) getCommonType(l, r llvm.Type) llvm.Type {
	if l == r {
		return l
	}

	if l == llvm.PointerType(r, 0) {
		return r
	}

	panic(fmt.Errorf("Abnormal hanling of types %s %s", l, r))
}

func (be BinaryExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	err, leftVal := be.Left.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}
	if leftVal == nil {
		return fmt.Errorf("left side of expression is nil"), nil
	}

	err, rightVal := be.Right.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}
	if rightVal == nil {
		return fmt.Errorf("right side of expression is nil"), nil
	}

	handler, ok := handlers[be.Operator.Kind]
	if !ok {
		panic(fmt.Errorf("unsupported operator <%s>", be.Operator.Kind))
	}

	ctype := commonType(leftVal.Type(), rightVal.Type())
	left := be.castToType(ctx, ctype, *leftVal)
	right := be.castToType(ctx, ctype, *rightVal)

	return handler(ctx, left, right)
}

func commonType(l, r llvm.Type) llvm.Type {
	if l == llvm.PointerType(r, 0) {
		return r
	}

	if l == r {
		return l
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
	res := ctx.Builder.CreateICmp(llvm.IntULE, l, r, "")
	return nil, &res
}
