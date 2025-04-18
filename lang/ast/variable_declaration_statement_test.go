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
	"testing"

	"github.com/llir/llvm/ir"
	"github.com/stretchr/testify/assert"
)

func TestVariableDeclarationConfig(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		stmt := VarDeclarationStatement{
			Name:  "v",
			Value: NumberExpression{Value: 1},
		}
		expectedLL := "@v = global i32 1\n"

		expectCompilesSuccessFully(t, stmt, expectedLL)
	})
	t.Run("string", func(t *testing.T) {
		stmt := VarDeclarationStatement{
			Name:  "s",
			Value: StringExpression{Value: "Na nga def"},
		}
		expectedLL := "@s = global [10 x i8] c\"Na nga def\"\n"

		expectCompilesSuccessFully(t, stmt, expectedLL)
	})
	t.Run("array of numbers", func(t *testing.T) {
		stmt := VarDeclarationStatement{
			Name: "s",
			Value: ArrayInitializationExpression{
				Underlying: SymbolType{Name: "number"},
				Contents: []Expression{
					NumberExpression{Value: 0},
					NumberExpression{Value: 1},
					NumberExpression{Value: 2},
					NumberExpression{Value: 3},
					NumberExpression{Value: 4},
				},
			},
		}
		expectedLL := "@s = global [5 x i32] [i32 0, i32 1, i32 2, i32 3, i32 4]\n"

		expectCompilesSuccessFully(t, stmt, expectedLL)
	})
}

func expectCompilesSuccessFully(t *testing.T, stmt VarDeclarationStatement, expectedLL string) {
	t.Helper()
	mod := ir.NewModule()
	ctx := NewContext(nil, mod)

	err := stmt.Compile(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedLL, mod.String())

	_, ok := ctx.vars[stmt.Name]
	assert.True(t, ok)
}
