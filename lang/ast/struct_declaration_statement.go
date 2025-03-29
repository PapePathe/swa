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
	"fmt"
	"swahili/lang/values"

	_types "github.com/llir/llvm/ir/types"
)

type StructProperty struct {
	PropType Type
}
type StructDeclarationStatement struct {
	Name       string
	Properties map[string]StructProperty
}

var _ Statement = (*StructDeclarationStatement)(nil)

func (cs StructDeclarationStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (ms StructDeclarationStatement) Compile(ctx *Context) error {
	attrs := []_types.Type{}

	for _, attr := range ms.Properties {
		switch v := attr.PropType.(type) {
		case SymbolType:
			switch v.Name {
			case "Chaine":
				attrs = append(attrs, _types.NewPointer(_types.I8))
			case "Nombre":
				attrs = append(attrs, _types.I32)
			default:
				err := fmt.Errorf("struct proprerty type %s not supported", v.Name)
				panic(err)
			}
		default:
			err := fmt.Errorf("struct proprerty does not support type")
			panic(err)
		}
	}

	ctx.mod.NewTypeDef(ms.Name, _types.NewStruct(attrs...))

	return nil
}
