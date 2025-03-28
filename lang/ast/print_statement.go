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
	"bytes"
	"encoding/json"
	"fmt"
	"swahili/lang/values"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type PrintStatetement struct {
	Values []Expression
}

var _ Statement = (*PrintStatetement)(nil)

func (ps PrintStatetement) Evaluate(s *Scope) (error, values.Value) {
	var buffer bytes.Buffer

	lg.Debug("Evaluating print statement", "stmt", ps)

	for _, v := range ps.Values {
		err, vals := v.Evaluate(s)
		if err != nil {
			return err, nil
		}

		buffer.WriteString(vals.String())
	}

	fmt.Println(buffer.String())

	return nil, values.StringValue{Value: buffer.String()}
}

func (cs PrintStatetement) statement() {}

func (ps PrintStatetement) Compile(ctx *Context) error {
	for _, v := range ps.Values {
		err, res := v.Compile(ctx)
		if err != nil {
			return err
		}

		switch v.(type) {
		case SymbolExpression:
			// gep := constant.NewGetElementPtr(
			// 	res.v.Type(),
			// 	res.c,
			// 	constant.NewInt(types.I32, 0),
			// 	constant.NewInt(types.I32, 0),
			// )
			ctx.NewCall(ir.NewFunc("printf", types.I32), res.v)
		case StringExpression:
			vl := ctx.NewAlloca(res.c.Type())
			ctx.NewStore(res.c, vl)

			ctx.NewCall(ir.NewFunc("printf", types.I32), vl)
		case NumberExpression:
			panic("NumberExpressio not implemented")
		default:
			err := fmt.Errorf("Print does not support <%s>", v)
			panic(err)
		}
	}

	return nil
}

func (cs PrintStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Values"] = cs.Values

	res := make(map[string]any)
	res["ast.PrintStatetement"] = m

	return json.Marshal(res)
}
