package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestSingleBranchConditional(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  ni(x>0) {
    r1.width = 100;
    r1.height = 100 + 400 - r1.width;
  }
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ConditionalStatetement{
				Condition: ast.BinaryExpression{
					Left: ast.SymbolExpression{Value: "x"},
					Right: ast.NumberExpression{
						Value: 0,
					}, Operator: lexer.Token{Value: ">", Kind: 18}},
				Success: ast.BlockStatement{
					Body: []ast.Statement{
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "=", Kind: 9},
								Assignee: ast.MemberExpression{
									Object:   ast.SymbolExpression{Value: "r1"},
									Property: ast.SymbolExpression{Value: "width"},
									Computed: false,
								},
								Value: ast.NumberExpression{Value: 100},
							}},
						ast.ExpressionStatement{
							Exp: ast.AssignmentExpression{
								Operator: lexer.Token{Value: "=", Kind: 9},
								Assignee: ast.MemberExpression{
									Object: ast.SymbolExpression{Value: "r1"},
									Property: ast.SymbolExpression{
										Value: "height",
									},
									Computed: false,
								},
								Value: ast.BinaryExpression{Left: ast.BinaryExpression{
									Left: ast.NumberExpression{
										Value: 100,
									}, Right: ast.NumberExpression{Value: 400},
									Operator: lexer.Token{Value: "+", Kind: 30}},
									Right: ast.MemberExpression{
										Object: ast.SymbolExpression{
											Value: "r1",
										}, Property: ast.SymbolExpression{Value: "width"}, Computed: false},
									Operator: lexer.Token{Value: "-", Kind: 22}},
							}}}},
				Failure: ast.BlockStatement{
					Body: []ast.Statement(nil),
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}
