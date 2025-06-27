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

	"tinygo.org/x/go-llvm"
)

type PrintStatetement struct {
	Values []Expression
}

var _ Statement = (*PrintStatetement)(nil)

func (ps PrintStatetement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	for _, v := range ps.Values {
		err, res := v.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}

		switch v.(type) {
		case StringExpression:
			all := ctx.Builder.CreateAlloca(res.Type(), "")
			ctx.Builder.CreateStore(*res, all)

			ctx.Builder.CreateCall(
				llvm.FunctionType(ctx.Context.Int32Type(), []llvm.Type{llvm.PointerType(ctx.Context.Int8Type(), 0)}, true),
				ctx.Module.NamedFunction("printf"),
				[]llvm.Value{all},
				"printf",
			)
		}

	}

	return nil, nil
}

func (cs PrintStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Values"] = cs.Values

	res := make(map[string]any)
	res["ast.PrintStatetement"] = m

	return json.Marshal(res)
}
