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

	"github.com/llir/llvm/ir/enum"
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

//func (be BinaryExpression) Evaluate(s *Scope) (error, values.Value) {
//	lg.Debug("Start", "node", be)
//
//	err, left := be.Left.Evaluate(s)
//	if err != nil {
//		lg.Error("ERROR evaluating left expression")
//
//		return err, nil
//	}
//
//	_, right := be.Right.Evaluate(s)
//
//	leftVal, _ := left.GetValue().(float64)
//	rightVal, _ := right.GetValue().(float64)
//
//	switch be.Operator.Kind {
//	case lexer.GreaterThan:
//		return nil, values.BooleaValue{Value: leftVal > rightVal}
//	case lexer.GreaterThanEquals:
//		return nil, values.BooleaValue{Value: leftVal >= rightVal}
//	case lexer.LessThan:
//		return nil, values.BooleaValue{Value: leftVal < rightVal}
//	case lexer.LessThanEquals:
//		return nil, values.BooleaValue{Value: leftVal <= rightVal}
//	default:
//		return fmt.Errorf("Operator not yet supportted %s", be.Operator.Kind), nil
//	}
//}

func (be BinaryExpression) Compile(ctx *Context) (error, *CompileResult) {
	err, leftVal := be.Left.Compile(ctx)
	if err != nil {
		return err, nil
	}
	err, rightVal := be.Right.Compile(ctx)
	if err != nil {
		return err, nil
	}

	switch be.Operator.Kind {
	case lexer.Plus:
		res := CompileResult{v: ctx.NewAdd(leftVal.c, rightVal.c)}
		return nil, &res
	case lexer.GreaterThan:
		res := CompileResult{v: ctx.NewICmp(enum.IPredUGT, leftVal.c, rightVal.c)}
		return nil, &res
	}

	return nil, nil
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
		return fmt.Errorf("left side of expression is nil"), nil
	}

	switch be.Operator.Kind {
	case lexer.Plus:
		res := ctx.Builder.CreateAdd(*leftVal, *rightVal, "")
		return nil, &res
	case lexer.GreaterThan:
		res := ctx.Builder.CreateICmp(llvm.IntUGT, *leftVal, *rightVal, "")
		return nil, &res
	case lexer.GreaterThanEquals:
		res := ctx.Builder.CreateICmp(llvm.IntUGE, *leftVal, *rightVal, "")
		return nil, &res
	case lexer.LessThan:
		res := ctx.Builder.CreateICmp(llvm.IntULT, *leftVal, *rightVal, "")
		return nil, &res
	case lexer.LessThanEquals:
		res := ctx.Builder.CreateICmp(llvm.IntULE, *leftVal, *rightVal, "")
		return nil, &res
	case lexer.Equals:
		res := ctx.Builder.CreateICmp(llvm.IntEQ, *leftVal, *rightVal, "")
		return nil, &res
	default:
		err := fmt.Errorf("unsupported operator <%s>", be.Operator.Kind)
		panic(err)
	}
}
