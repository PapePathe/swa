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
	"swahili/lang/values"
)

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}
func (cs ConditionalStatetement) statement() {}

func (cs ConditionalStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["success"] = cs.Success
	m["condition"] = cs.Condition
	m["failure"] = cs.Failure

	res := make(map[string]any)
	res["ast.ConditionalStatetement"] = m

	return json.Marshal(res)
}
