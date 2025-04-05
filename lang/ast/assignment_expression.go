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
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

// AssignmentExpression.
// Is an expression where the programmer is trying to assign a value to a variable.
//
// a = a +5;
// foo.bar = foo.bar + 10;
type AssignmentExpression struct {
	Operator lexer.Token
	Assignee Expression
	Value    Expression
}

var _ Expression = (*AssignmentExpression)(nil)

func (ae AssignmentExpression) Compile(ctx *Context) (error, *CompileResult) {
	err, receiver := ae.Assignee.Compile(ctx)
	if err != nil {
		return err, nil
	}

	err, value := ae.Value.Compile(ctx)
	if err != nil {
		return err, nil
	}

	alloc := ctx.NewAlloca(receiver.c.Type())
	ctx.NewStore(value.c, alloc)

	return nil, &CompileResult{v: alloc}
}

func (bs AssignmentExpression) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	return nil, nil
}
