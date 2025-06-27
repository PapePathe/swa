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

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body []Statement
}

var _ Statement = (*BlockStatement)(nil)

func (bs BlockStatement) Compile(ctx *Context) error {
	for _, stmt := range bs.Body {
		err := stmt.Compile(ctx)
		if err != nil {
			lg.Error("ERROR", " evaluating statement", err.Error())

			return err
		}
	}

	return nil
}

func (bs BlockStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	for _, stmt := range bs.Body {
		err, _ := stmt.CompileLLVM(ctx)
		if err != nil {
			lg.Error("ERROR", " evaluating statement", err.Error())

			return err, nil
		}
	}

	return nil, nil
}

func (bs BlockStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["ast.BlockStatement"] = bs.Body

	return json.Marshal(m)
}
