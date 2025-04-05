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

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"tinygo.org/x/go-llvm"
)

// NumberExpression ...
type NumberExpression struct {
	Value float64
}

var _ Expression = (*NumberExpression)(nil)

func (nexpr NumberExpression) Compile(ctx *Context) (error, *CompileResult) {
	return nil, &CompileResult{c: constant.NewInt(types.I32, int64(nexpr.Value))}
}

func (cs NumberExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["value"] = cs.Value

	res := make(map[string]any)
	res["ast.NumberExpression"] = m

	return json.Marshal(res)
}

func (se NumberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), false)

	return nil, &res
}
