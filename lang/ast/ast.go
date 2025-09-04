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
	"tinygo.org/x/go-llvm"
)

// Statement ...
type Statement interface {
	//	Compile(ctx *Context) error
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

// Expression ...
type Expression interface {
	//	Compile(ctx *Context) (error, *CompileResult)
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

// Type
type Type interface {
	_type()
}

type CompilerCtx struct {
	Context     *llvm.Context
	Builder     *llvm.Builder
	Module      *llvm.Module
	SymbolTable map[string]llvm.Value
}
