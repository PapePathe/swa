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
	"github.com/llir/llvm/ir/constant"
	"tinygo.org/x/go-llvm"
)

// StringExpression ...
type StringExpression struct {
	Value string
}

var _ Expression = (*StringExpression)(nil)

func (se StringExpression) Compile(ctx *Context) (error, *CompileResult) {
	return nil, &CompileResult{c: constant.NewCharArrayFromString(se.Value)}
}

func (se StringExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	res := ctx.Context.ConstString(se.Value, true)
	return nil, &res
}
