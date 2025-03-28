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
	"swahili/lang/lexer"
)

type PrefixExpression struct {
	Operator        lexer.Token
	RightExpression Expression
}

var _ Expression = (*PrefixExpression)(nil)

func (cs PrefixExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["operator"] = cs.Operator
	m["right_expression"] = cs.RightExpression

	res := make(map[string]any)
	res["ast.PrefixExpression"] = m

	return json.Marshal(res)
}

func (PrefixExpression) Compile(ctx *Context) (error, *CompileResult) {
	return nil, nil
}
