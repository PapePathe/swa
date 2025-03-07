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

package parser_test

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedASTForArrayDecl = ast.BlockStatement{
	Body: []ast.Statement{
		ast.VarDeclarationStatement{
			Name:       "numbers",
			IsConstant: true,
			Value: ast.ArrayInitializationExpression{
				Underlying: ast.SymbolType{Name: "number"},
				Contents: []ast.Expression{
					ast.NumberExpression{Value: 1},
					ast.NumberExpression{Value: 2},
					ast.NumberExpression{Value: 3},
					ast.NumberExpression{Value: 4},
					ast.NumberExpression{Value: 5},
					ast.NumberExpression{Value: 6},
				},
			},
			ExplicitType: ast.Type(nil),
		},
	},
}

func TestArrayInstantiationEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:english; 
			const numbers = []number{1,2,3,4,5,6};
	`))
	assert.Equal(t, result, expectedASTForArrayDecl)
}

func TestArrayInstantiationFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:french; 
			constante numbers = []number{1,2,3,4,5,6};
	`))

	assert.Equal(t, result, expectedASTForArrayDecl)
}

func TestArrayInstantiationMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke; 
			const numbers = []number{1,2,3,4,5,6};
	`))

	assert.Equal(t, result, expectedASTForArrayDecl)
}
