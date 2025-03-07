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

var expectedASTForDualCondition = ast.BlockStatement{
	Body: []ast.Statement{
		ast.ConditionalStatetement{
			Condition: ast.BinaryExpression{
				Left:     ast.SymbolExpression{Value: "x"},
				Right:    ast.NumberExpression{Value: 0},
				Operator: lexer.Token{Value: ">", Kind: 18, Name: "GREATER_THAN"},
			},
			Success: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
							Assignee: ast.SymbolExpression{Value: "width"},
							Value:    ast.NumberExpression{Value: 100},
						},
					},
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
							Assignee: ast.SymbolExpression{Value: "height"},
							Value: ast.BinaryExpression{
								Left: ast.BinaryExpression{
									Left:     ast.NumberExpression{Value: 100},
									Right:    ast.NumberExpression{Value: 400},
									Operator: lexer.Token{Value: "+", Kind: 30, Name: "PLUS"},
								},
								Right:    ast.SymbolExpression{Value: "width"},
								Operator: lexer.Token{Value: "-", Kind: 22, Name: "MINUS"},
							},
						},
					},
				},
			},
			Failure: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "+=", Kind: 31, Name: "PLUS_EQUAL"},
							Assignee: ast.SymbolExpression{Value: "width"},
							Value:    ast.NumberExpression{Value: 300},
						},
					},
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "+=", Kind: 31, Name: "PLUS_EQUAL"},
							Assignee: ast.SymbolExpression{Value: "height"},
							Value: ast.BinaryExpression{
								Left:     ast.SymbolExpression{Value: "rand"},
								Right:    ast.NumberExpression{Value: 14},
								Operator: lexer.Token{Value: "/", Kind: 15, Name: "DIVIDE"},
							},
						},
					},
				},
			},
		},
	},
}

func TestDualBranchConditionalFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  	dialect:french;
    si(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } sinon {
      width += 300;
      height += rand / 14;
    }
		`))

	assert.Equal(t, expectedASTForDualCondition, result)
}

func TestDualBranchConditionalMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  	dialect:malinke;
    ni(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } nii {
    	width += 300;
    	height += rand / 14;
    }
	`))

	assert.Equal(t, expectedASTForDualCondition, result)
}

func TestDualBranchConditionalEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  	dialect:english;
    if(x>0) {
      width = 100;
      height = 100 + 400 - width;
    } else {
    	width += 300;
    	height += rand / 14;
    }
	`))

	assert.Equal(t, expectedASTForDualCondition, result)
}

var expectedAST = ast.BlockStatement{
	Body: []ast.Statement{
		ast.ConditionalStatetement{
			Condition: ast.BinaryExpression{
				Left: ast.SymbolExpression{Value: "x"},
				Right: ast.NumberExpression{
					Value: 0,
				}, Operator: lexer.Token{Value: ">", Kind: 18, Name: "GREATER_THAN"},
			},
			Success: ast.BlockStatement{
				Body: []ast.Statement{
					ast.ExpressionStatement{
						Exp: ast.AssignmentExpression{
							Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
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
							Operator: lexer.Token{Value: "=", Kind: 9, Name: "ASSIGNMENT"},
							Assignee: ast.MemberExpression{
								Object: ast.SymbolExpression{Value: "r1"},
								Property: ast.SymbolExpression{
									Value: "height",
								},
								Computed: false,
							},
							Value: ast.BinaryExpression{
								Left: ast.BinaryExpression{
									Left: ast.NumberExpression{
										Value: 100,
									}, Right: ast.NumberExpression{Value: 400},
									Operator: lexer.Token{Value: "+", Kind: 30, Name: "PLUS"},
								},
								Right: ast.MemberExpression{
									Object: ast.SymbolExpression{
										Value: "r1",
									}, Property: ast.SymbolExpression{Value: "width"}, Computed: false,
								},
								Operator: lexer.Token{Value: "-", Kind: 22, Name: "MINUS"},
							},
						},
					},
				},
			},
			Failure: ast.BlockStatement{
				Body: []ast.Statement(nil),
			},
		},
	},
}

func TestSingleBranchConditionalFrench(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:french;
 		 si(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}

func TestSingleBranchConditionalMalinke(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:malinke;
 		 ni(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}

func TestSingleBranchConditionalEnglish(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
 		 dialect:english;
 		 if(x>0) {
 		   r1.width = 100;
 		   r1.height = 100 + 400 - r1.width;
 		 }
	`))

	assert.Equal(t, expectedAST, result)
}
