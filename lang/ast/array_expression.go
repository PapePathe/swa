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

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (v ArrayInitializationExpression) Compile(ctx *Context) (error, *CompileResult) {
	contents := []constant.Constant{}

	for _, elem := range v.Contents {
		err, c := elem.Compile(ctx)
		if err != nil {
			return err, nil
		}

		contents = append(contents, c.c)
	}

	return nil, &CompileResult{c: constant.NewArray(nil, contents...)}
}

func (bs ArrayInitializationExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
