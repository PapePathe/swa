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
	"fmt"
	"swahili/lang/values"
)

// ExpressionStatement ...
type ExpressionStatement struct {
	// The expression
	Exp Expression
}

var _ Statement = (*ExpressionStatement)(nil)

func (cs ExpressionStatement) Evaluate(s *Scope) (error, values.Value) {
	fmt.Println("Evaluating expression statement")
	return nil, nil
}

func (bs ExpressionStatement) statement() {}

func (cs ExpressionStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["expression"] = cs.Exp

	res := make(map[string]any)
	res["ast.ExpressionStatetement"] = m

	return json.Marshal(res)
}
