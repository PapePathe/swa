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

import "tinygo.org/x/go-llvm"

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

var _ Expression = (*MemberExpression)(nil)

func (MemberExpression) Compile(ctx *Context) (error, *CompileResult) {
	return nil, nil
}

func (bs MemberExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
