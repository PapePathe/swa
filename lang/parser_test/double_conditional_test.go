package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestDualBranchConditional(t *testing.T) {
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

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ConditionalStatetement{
				Condition: ast.BinaryExpression{
					Left:     ast.SymbolExpression{Value: "x"},
					Right:    ast.NumberExpression{Value: 0},
					Operator: lexer.Token{Value: ">", Kind: 18},
				},
				Success: ast.BlockStatement{
					Body: []ast.Statement{
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "=", Kind: 9},
								Assignee: ast.SymbolExpression{Value: "width"},
								Value:    ast.NumberExpression{Value: 100},
							},
						},
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "=", Kind: 9},
								Assignee: ast.SymbolExpression{Value: "height"},
								Value: ast.BinaryExpression{
									Left: ast.BinaryExpression{
										Left:     ast.NumberExpression{Value: 100},
										Right:    ast.NumberExpression{Value: 400},
										Operator: lexer.Token{Value: "+", Kind: 30},
									},
									Right:    ast.SymbolExpression{Value: "width"},
									Operator: lexer.Token{Value: "-", Kind: 22},
								},
							},
						},
					},
				},
				Failure: ast.BlockStatement{
					Body: []ast.Statement{
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "+=", Kind: 31},
								Assignee: ast.SymbolExpression{Value: "width"},
								Value:    ast.NumberExpression{Value: 300},
							},
						},
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "+=", Kind: 31},
								Assignee: ast.SymbolExpression{Value: "height"},
								Value: ast.BinaryExpression{
									Left:     ast.SymbolExpression{Value: "rand"},
									Right:    ast.NumberExpression{Value: 14},
									Operator: lexer.Token{Value: "/", Kind: 15},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}
