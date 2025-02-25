package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/parser"
)

func TestStructPropertyAssignment(t *testing.T) {
	result := parser.Parse(lexer.Tokenize(`
  		dialect:  malinke; 
      r1.width += 100;
      r1.width += rand / 14;
	`))

	expected := ast.BlockStatement{
		Body: []ast.Statement{
			ast.ExpressionStatement{
				Exp: ast.AssignmentExpression{
					Operator: lexer.Token{Value: "+=", Kind: 31},
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
					Operator: lexer.Token{Value: "+=", Kind: 31},
					Assignee: ast.MemberExpression{
						Object:   ast.SymbolExpression{Value: "r1"},
						Property: ast.SymbolExpression{Value: "width"},
						Computed: false,
					},
					Value: ast.BinaryExpression{
						Left:     ast.SymbolExpression{Value: "rand"},
						Right:    ast.NumberExpression{Value: 14},
						Operator: lexer.Token{Value: "/", Kind: 15},
					},
				},
			},
		},
	}

	assert.Equal(t, result, expected)
}
