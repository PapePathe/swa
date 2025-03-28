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

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (v VarDeclarationStatement) Evaluate(s *Scope) (error, values.Value) {
	lg.Debug("Evaluating variable declaration statement", "variable", v)

	err, val := v.Value.Evaluate(s)
	if err != nil {
		return err, nil
	}

	lg.Debug("Result of evaluating value", "value", val)

	s.Set(v.Name, val)

	return nil, nil
}

func (bs VarDeclarationStatement) statement() {}

func (vd VarDeclarationStatement) Compile(ctx *Context) error {
	err, cst := vd.Value.Compile(ctx)

	if err != nil {
		return err
	}

	if ctx.parent == nil {
		ctx.vars[vd.Name] = Var{
			cst: cst.c,
			def: ctx.mod.NewGlobalDef(vd.Name, cst.c),
		}
	} else {
		// TODO handle case where variable is local to the current block
		// storage := ctx.NewAlloca(cst.c.Type())
		// ctx.NewStore(cst.c, storage)
	}

	return nil
}

func (cs VarDeclarationStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["name"] = cs.Name
	m["is_constant"] = cs.IsConstant
	m["value"] = cs.Value
	m["explicit_type"] = cs.ExplicitType

	res := make(map[string]any)
	res["ast.VarDeclarationStatement"] = m

	return json.Marshal(res)
}
