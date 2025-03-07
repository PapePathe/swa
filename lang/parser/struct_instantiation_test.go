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

func TestStructInstantiation(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:malinke; 
      const r1: Rectangle = Rectangle {
        width: 10.2,
        height: 45.2 + 6 - 11 * 999,
      };
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.VarDeclarationStatement{
				Name:       "r1",
				IsConstant: true,
				Value: ast.StructInitializationExpression{
					Name: "Rectangle",
					Properties: map[string]ast.Expression{
						"height": ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:     ast.NumberExpression{Value: 45.2},
								Right:    ast.NumberExpression{Value: 6},
								Operator: lexer.Token{Value: "+", Kind: 30, Name: "PLUS"},
							},
							Right: ast.BinaryExpression{
								Left:     ast.NumberExpression{Value: 11},
								Right:    ast.NumberExpression{Value: 999},
								Operator: lexer.Token{Value: "*", Kind: 32, Name: "STAR"},
							},
							Operator: lexer.Token{Value: "-", Kind: 22, Name: "MINUS"},
						},
						"width": ast.NumberExpression{Value: 10.2},
					},
				},
				ExplicitType: ast.SymbolType{Name: "Rectangle"},
			},
		},
	}

	assert.Equal(t, result, expected)
}

var expectedAstForTestStructDeclaration = ast.BlockStatement{
	Body: []ast.Statement{
		ast.StructDeclarationStatement{
			Name: "Rectangle",
			Properties: map[string]ast.StructProperty{
				"height": {PropType: ast.SymbolType{Name: "float"}},
				"name":   {PropType: ast.SymbolType{Name: "string"}},
				"width":  {PropType: ast.SymbolType{Name: "number"}},
			},
		},
	},
}

func TestStructDeclarationFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:french;
      structure Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))
	assert.Equal(t, result, expectedAstForTestStructDeclaration)
}

func TestStructDeclarationMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:malinke;
      struct Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))
	assert.Equal(t, result, expectedAstForTestStructDeclaration)
}

func TestStructDeclarationEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
      dialect:english;
      struct Rectangle {
        width: number,
        height: float,
        name: string,
      }
	`))
	assert.Equal(t, result, expectedAstForTestStructDeclaration)
}

func TestStructPropertyAssignmentFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:malinke; 
      r1.width += 100;
      r1.width += rand / 14;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "+=", Kind: 31, Name: "PLUS_EQUAL"},
					Assignee: ast.MemberExpression{
						Object:   ast.SymbolExpression{Value: "r1"},
						Property: ast.SymbolExpression{Value: "width"},
						Computed: false,
					},
					Value: ast.NumberExpression{Value: 100},
				},
			},
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "+=", Kind: 31, Name: "PLUS_EQUAL"},
					Assignee: ast.MemberExpression{
						Object:   ast.SymbolExpression{Value: "r1"},
						Property: ast.SymbolExpression{Value: "width"},
						Computed: false,
					},
					Value: ast.BinaryExpression{
						Left:     ast.SymbolExpression{Value: "rand"},
						Right:    ast.NumberExpression{Value: 14},
						Operator: lexer.Token{Value: "/", Kind: 15, Name: "DIVIDE"},
					},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}
