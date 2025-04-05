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

	"tinygo.org/x/go-llvm"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value string
}

var _ Expression = (*SymbolExpression)(nil)

func (se SymbolExpression) Compile(ctx *Context) (error, *CompileResult) {
	val, err := ctx.LookupVariable(se.Value)
	if err != nil {
		localVar, err := ctx.LookupLocalVariable(se.Value)
		if err != nil {
			return err, nil
		}
		l := ctx.NewLoad(localVar.Value.Typ, localVar.Value)

		return nil, &CompileResult{v: l.Src}
	}

	return nil, &CompileResult{v: val.def, c: val.cst}
}

func (se SymbolExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	val, ok := ctx.SymbolTable[se.Value]

	if !ok {
		return fmt.Errorf("Variable %s does not exist", se.Value), nil
	}

	return nil, &val
}
