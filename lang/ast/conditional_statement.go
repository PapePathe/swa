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
	"encoding/json"
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"tinygo.org/x/go-llvm"
)

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) Compile(ctx *Context) error {
	successBlk := ctx.Parent.NewBlock("if.success")
	failureBlk := ctx.Parent.NewBlock("if.failure")
	//	mergeBlk := ctx.Parent.NewBlock("")
	//	successBlk.NewBr(mergeBlk)

	err, res := cs.Condition.Compile(ctx)
	// TODO fix issue with conditional branch
	if err != nil {
		return err
	}

	successCtx := ctx.NewContext(successBlk)
	successCtx.NewRet(constant.NewInt(types.I32, 0))

	err = cs.Success.Compile(successCtx)
	if err != nil {
		return err
	}

	failureCtx := ctx.NewContext(failureBlk)

	failureCtx.NewRet(constant.NewInt(types.I32, 0))

	err = cs.Failure.Compile(failureCtx)
	if err != nil {
		return err
	}

	ctx.NewCondBr(res.c, successBlk, failureBlk)
	return nil
}

func (cs ConditionalStatetement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	err, condition := cs.Condition.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}
	fmt.Println(condition.String())

	mergeBlock := ctx.Context.InsertBasicBlock(ctx.Builder.GetInsertBlock(), "merge")
	thenBlock := ctx.Context.InsertBasicBlock(ctx.Builder.GetInsertBlock(), "if")
	elseBlock := ctx.Context.InsertBasicBlock(ctx.Builder.GetInsertBlock(), "else")

	ctx.Builder.CreateCondBr(*condition, thenBlock, elseBlock)

	ctx.Builder.SetInsertPointAtEnd(thenBlock)
	err, successVal := cs.Success.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}
	ctx.Builder.CreateBr(mergeBlock)

	ctx.Builder.SetInsertPointAtEnd(elseBlock)
	err, failureVal := cs.Failure.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}
	ctx.Builder.CreateBr(mergeBlock)
	ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	var phi llvm.Value
	if successVal != nil {
		phi := ctx.Builder.CreatePHI(successVal.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal}, []llvm.BasicBlock{thenBlock})
	}

	if failureVal != nil {
		phi := ctx.Builder.CreatePHI(successVal.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal}, []llvm.BasicBlock{thenBlock})
	}
	return nil, &phi
}

func (cs ConditionalStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["success"] = cs.Success
	m["condition"] = cs.Condition
	m["failure"] = cs.Failure

	res := make(map[string]any)
	res["ast.ConditionalStatetement"] = m

	return json.Marshal(res)
}
